package main

import (
	"../bitoip"
	"github.com/golang/glog"
	"net"
	"strings"
	"time"
)
/**
	Handle an incoming message to the reflector
 */
func Handler(serverAddress *net.UDPAddr, msg bitoip.RxMSG) {
	switch msg.Verb {
	case bitoip.EnumerateChannels:
		responsePayload := new(bitoip.ListChannelsPayload)
		copy(responsePayload.Channels[:], ChannelIds())
		bitoip.UDPTx(bitoip.ListChannels,
					 msg.Payload,
					 &msg.SrcAddress)

	case bitoip.CarrierEvent:
		ce := msg.Payload.(*bitoip.CarrierEventPayload)
		glog.V(1).Infof("got carrier event %v", ce)
		channel := GetChannel(ce.Channel)
		channel.Subscribe(msg.SrcAddress, "????????") //make sure this user subscribed
		channel.Broadcast(*ce)

	case bitoip.ListenRequest:
		lr := msg.Payload.(*bitoip.ListenRequestPayload)
		channel := GetChannel(lr.Channel)
		key := channel.Subscribe(msg.SrcAddress, strings.Trim(string(lr.Callsign[:]), "\x00"))
		lcp := bitoip.ListenConfirmPayload{lr.Channel, key}

		bitoip.UDPTx(bitoip.ListenConfirm, lcp, &msg.SrcAddress)

	case bitoip.TimeSync:
		ts := msg.Payload.(*bitoip.TimeSyncPayload)

		tsr := bitoip.TimeSyncResponsePayload{
			ts.CurrentTime,
		   time.Now().UnixNano(),
		}

		bitoip.UDPTx(bitoip.TimeSyncResponse, tsr, &msg.SrcAddress)
	}

}
