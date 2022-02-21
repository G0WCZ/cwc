/*
Copyright 2022 Graeme Sutherland


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"context"
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/G0WCZ/cwc/cwcpb"
	"net"
	"sort"
	"testing"
	"time"

	"github.com/golang/glog"
	"gotest.tools/assert"
)

func TestNewChannel(t *testing.T) {
	channel := NewChannel(33)
	assert.DeepEqual(t, channel.ChannelId, uint32(33))
	assert.Equal(t, len(channel.Subscribers), 0)
	assert.Equal(t, len(channel.Addresses), 0)
	assert.Equal(t, channel.LastKey, uint32(99))
}

func TestGetChannel(t *testing.T) {
	channel1 := GetChannel(21)
	channel2 := GetChannel(21)
	assert.DeepEqual(t, channel1, channel2)
}

func TestSubscribeWhenNotSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}

func TestSubscribeWhenSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
	channel1.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}
func TestUnsubscribeWhenSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel2.Addresses), 1)
	assert.Equal(t, len(channel2.Subscribers), 1)
	channel2.Unsubscribe(*addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func TestUnsubscribeWhenNotSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Unsubscribe(*addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func sortSlice(sl []uint32) []uint32 {
	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })
	return sl
}

func TestChannelIds(t *testing.T) {
	GetChannel(21)
	GetChannel(22)
	GetChannel(33)
	assert.DeepEqual(t, sortSlice(ChannelIds()), sortSlice([]uint32{21, 22, 33}))
}

func TestEmptyChannelIds(t *testing.T) {
	channels = make(map[uint32]*Channel)
	assert.Equal(t, len(ChannelIds()), 0)
}

func carrierEventPayload(key bitoip.CarrierKeyType) *cwcpb.CWCMessage {
	beOn := &cwcpb.CarrierEvent_BitEvent{
		BitEvent:   true,
		TimeOffset: 0,
		Last:       false,
	}

	beOff := &cwcpb.CarrierEvent_BitEvent{
		BitEvent:   false,
		TimeOffset: 100,
		Last:       true,
	}
	return &cwcpb.CWCMessage{
		Msg: &cwcpb.CWCMessage_CarrierEvent{
			CarrierEvent: &cwcpb.CarrierEvent{
				ChannelId:      1,
				CarrierKey:     key,
				StartTimestamp: time.Now().UnixNano(),
				SendTimestamp:  time.Now().UnixNano(),
				BitEvents:      []*cwcpb.CarrierEvent_BitEvent{beOn, beOff},
			},
		},
	}
}

func TestBroadcastToSubscriber(t *testing.T) {
	c1 := GetChannel(1)
	add := "localhost:9453"
	addr, _ := net.ResolveUDPAddr("udp", add)
	glog.Infof("addr: %v", addr)
	key := c1.Subscribe(*addr, "G0WCZ")
	ce := carrierEventPayload(key)

	pc, _ := net.ListenPacket("udp", add)
	buffer := make([]byte, bitoip.MaxMessageSizeInBytes)
	doneChan := make(chan []byte, 1)

	// get one message
	go func() {
		n, _, _ := pc.ReadFrom(buffer)
		glog.Infof("Raw Rx: %d %v", n, buffer)
		doneChan <- buffer[0:n]
	}()

	serverAddress, _ := net.ResolveUDPAddr("udp", "localhost:6012")
	ctx := context.Background()
	messages := make(chan bitoip.RxMSG)
	go bitoip.UDPRx(ctx, serverAddress, messages)

	// delay for connection to be established
	time.Sleep(time.Second * 2)

	// broadcast
	c1.Broadcast(ce)

	buf := <-doneChan

	msg := bitoip.DecodeBuffer(buf, len(buf))
	assert.Equal(t, msg.GetCarrierEvent().GetChannelId(), uint32(1))
}

func TestSuperviseChannelsNoSubscribers(t *testing.T) {
	channels = make(map[uint32]*Channel)
	_ = GetChannel(1)
	_ = GetChannel(2)
	r := SuperviseChannels(time.Now(), time.Duration(10*time.Minute))
	assert.Equal(t, r, 0)
}

func TestSuperviseChannelsNoneRemoved(t *testing.T) {
	channels = make(map[uint32]*Channel)
	c1 := GetChannel(1)
	c2 := GetChannel(2)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:19234")
	c1.Subscribe(*addr, "A1AAA")
	c2.Subscribe(*addr, "A1BBB")
	r := SuperviseChannels(time.Now(), time.Duration(10*time.Minute))
	assert.Equal(t, r, 0)
}

func TestSuperviseChannels2Removed(t *testing.T) {
	channels = make(map[uint32]*Channel)
	c1 := GetChannel(1)
	c2 := GetChannel(2)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:19234")
	c1.Subscribe(*addr, "A1AAA")
	c2.Subscribe(*addr, "A1BBB")
	r := SuperviseChannels(time.Now().Add(time.Duration(20*time.Minute)), time.Duration(10*time.Minute))
	assert.Equal(t, r, 2)
	r = SuperviseChannels(time.Now().Add(time.Duration(20*time.Minute)), time.Duration(10*time.Minute))
	assert.Equal(t, r, 0)
}

func TestStation_AddSeenOn(t *testing.T) {
	s := Station{
		"A1AAA",
		nil,
		time.Now(),
		1,
		time.Now(),
	}
	s.AddSeenOn(bitoip.ChannelIdType(1))
	assert.Equal(t, uint32(1), s.SeenOnChannels[0])
	s.AddSeenOn(bitoip.ChannelIdType(1))
	assert.Equal(t, 1, len(s.SeenOnChannels))
	s.AddSeenOn(bitoip.ChannelIdType(2))
	s.AddSeenOn(bitoip.ChannelIdType(3))
	s.AddSeenOn(bitoip.ChannelIdType(4))
	s.AddSeenOn(bitoip.ChannelIdType(5))
	assert.DeepEqual(t,
		s.SeenOnChannels,
		[]bitoip.ChannelIdType{5, 4, 3, 2, 1},
	)
}
