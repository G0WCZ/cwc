package main

import (
	"net"
	"sort"
	"testing"
	"gotest.tools/assert"
	"../bitoip"
	"time"
)

func TestNewChannel(t *testing.T) {
	channel := NewChannel(33)
	assert.DeepEqual(t,channel.ChannelId, uint16(33))
	assert.Equal(t, len(channel.Subscribers), 0)
	assert.Equal(t, len(channel.Addresses), 0)
}

func TestGetChannel(t *testing.T) {
	channel1 := GetChannel(21)
	channel2 := GetChannel(21)
	assert.DeepEqual(t, channel1, channel2)
}

func TestSubscribeWhenNotSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(addr)
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}

func TestSubscribeWhenSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(addr)
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
	channel1.Subscribe(addr)
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}
func TestUnsubscribeWhenSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Subscribe(addr)
	assert.Equal(t, len(channel2.Addresses), 1)
	assert.Equal(t, len(channel2.Subscribers), 1)
	channel2.Unsubscribe(addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func TestUnsubscribeWhenNotSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Unsubscribe(addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func sortSlice(sl []uint16) []uint16 {
	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })
	return sl
}

func TestChannelIds(t *testing.T) {
	GetChannel(21)
	GetChannel(22)
	GetChannel(33)
	assert.DeepEqual(t, sortSlice(ChannelIds()), sortSlice([]uint16{21, 22, 33}))
}

func TestEmptyChannelIds(t *testing.T) {
	channels = make(map[uint16]Channel)
	assert.Equal(t, len(ChannelIds()), 0)
}

func carrierEventPayload() bitoip.CarrierEventPayload {
	return bitoip.CarrierEventPayload{
		1,
		99,
		time.Now().Unix(),
		[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{
			bitoip.CarrierBitEvent{0, bitoip.BitOn},
			bitoip.CarrierBitEvent{100, bitoip.BitOff & bitoip.LastEvent},
		},
	}
}

func TestBroadcastEmpty(t *testing.T) {
	channels = make(map[uint16]Channel)
	c1 := GetChannel(1)
	ce := carrierEventPayload()

	c1.Broadcast(ce, nil)
}

func TestBroadcastToSubscriber(t *testing.T) {
	channels = make(map[uint16]Channel)
	c1 := GetChannel(1)
	ce := carrierEventPayload()
	add := "localhost:2020"
	addr, _ := net.ResolveUDPAddr("udp", add)

	c1.Subscribe(addr)


	pc, _ := net.ListenPacket("udp", add)
	buffer := make([]byte, bitoip.MaxMessageSizeInBytes)
	doneChan := make(chan []byte, 1)

	// get one message
	go func() {
		_, _, _ = pc.ReadFrom(buffer)
		doneChan <- buffer
	}()

	// broadcast
	c1.Broadcast(ce, nil)

	buf :=  <- doneChan

	verb, payload := bitoip.DecodePacket(buf)
	assert.Equal(t, verb, bitoip.CarrierEvent)
	assert.DeepEqual(t, payload, &ce)
	rxce := payload.(*bitoip.CarrierEventPayload)
	assert.Equal(t, rxce.CarrierKey, uint16(99))
}