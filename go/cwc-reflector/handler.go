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
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/G0WCZ/cwc/cwcpb"
	"net"
	"time"

	"github.com/golang/glog"
)

// Handle an incoming message to the reflector
func Handler(serverAddress *net.UDPAddr, msg bitoip.RxMSG) {
	switch msg.Message.Msg.(type) {
	// Channel list
	case *cwcpb.CWCMessage_EnumerateChannels:
		cwcmsg := &cwcpb.CWCMessage{}
		lc := &cwcpb.CWCMessage_ListChannels{}
		lc.ListChannels.ChannelIds = append(lc.ListChannels.ChannelIds, ChannelIds()...)
		cwcmsg.Msg = lc
		bitoip.UDPTx(cwcmsg, &msg.SrcAddress)
	// Carrier morse data
	case *cwcpb.CWCMessage_CarrierEvent:
		msg := msg.Message
		glog.V(1).Infof("got carrier event %v", msg)
		channel := GetChannel(msg.GetCarrierEvent().GetChannelId())
		channel.Broadcast(msg)

	// Subscribe request
	case *cwcpb.CWCMessage_ListenRequest:
		lr := msg.Message.GetListenRequest()
		channel := GetChannel(lr.GetChannelId())
		key := channel.Subscribe(msg.SrcAddress, lr.Callsign)
		lcp := &cwcpb.CWCMessage_ListenConfirm{
			ListenConfirm: &cwcpb.ListenConfirm{
				ChannelId:  lr.GetChannelId(),
				CarrierKey: key,
			},
		}
		lcMsg := &cwcpb.CWCMessage{
			Msg: lcp,
		}
		bitoip.UDPTx(lcMsg, &msg.SrcAddress)

	// Time sync
	case *cwcpb.CWCMessage_TimeSync:
		tsr := msg.Message.GetTimeSync()
		response := &cwcpb.CWCMessage_TimeSyncResponse{
			TimeSyncResponse: &cwcpb.TimeSyncResponse{
				GivenTime:    tsr.GetCurrentTime(),
				ServerRxTime: msg.RxTime,
				ServerTxTime: time.Now().UnixNano(),
			},
		}
		responseMsg := &cwcpb.CWCMessage{
			Msg: response,
		}
		bitoip.UDPTx(responseMsg, &msg.SrcAddress)
	}
}
