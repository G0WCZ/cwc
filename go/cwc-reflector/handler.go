package main

import (
	"../bitoip"
	"log"
	"net"
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
		log.Printf("got carrier event %v", ce)
		channel := GetChannel(ce.Channel)
		channel.Subscribe(msg.SrcAddress) //make sure this user subscribed
		channel.Broadcast(*ce)

	case bitoip.ListenRequest:
		lr := msg.Payload.(*bitoip.ListenRequestPayload)
		channel := GetChannel(lr.Channel)
		key := channel.Subscribe(msg.SrcAddress)
		lcp := bitoip.ListenConfirmPayload{lr.Channel, key}

		bitoip.UDPTx(bitoip.ListenConfirm, lcp, &msg.SrcAddress)
	}
}
