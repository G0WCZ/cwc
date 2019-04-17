package main

import (
	"flag"
	"log"
	"net"
			"context"
	"../bitoip"
	"os"
)

var serverAddress *net.UDPAddr

func main() {
	address := flag.String("address", "localhost:7388", "-address=host:port")

	flag.Parse()

	ReflectorServer(context.TODO(), *address)
}

func ReflectorServer(ctx context.Context, address string) {

	serverAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Can't use address %s: %s", address, err)
		os.Exit(1)
	}
	log.Printf("Starting reflector on %s", address)

	messages := make(chan bitoip.RxMSG)

	bitoip.UDPRx(ctx, *serverAddress, messages)

	for {
		select {
		case <- ctx.Done():
			return
		case m := <- messages:
			Handler(m)
		}
	}
}

/**
	Handle an incoming message to the reflector
 */
func Handler(msg bitoip.RxMSG) {
	switch msg.Verb {
	case bitoip.EnumerateChannels:
		responsePayload := new(bitoip.ListChannelsPayload)
		copy(responsePayload.Channels[:], ChannelIds())
		bitoip.UDPTx(bitoip.ListChannels,
					 msg.Payload,
					 msg.SrcAddress.String(),
					 serverAddress)

	case bitoip.CarrierEvent:
		ce := msg.Payload.(bitoip.CarrierEventPayload)
		channel := GetChannel(ce.Channel)
		channel.Subscribe(msg.SrcAddress) //make sure this user subscribed
		channel.Broadcast(ce, serverAddress)

	case bitoip.ListenRequest:
		lr := msg.Payload.(bitoip.ListenRequestPayload)
		channel := GetChannel(lr.Channel)
		key := channel.Subscribe(msg.SrcAddress)
		lcp := bitoip.ListenConfirmPayload{lr.Channel, key}

		bitoip.UDPTx(bitoip.ListenConfirm, lcp, msg.SrcAddress.String(), serverAddress)
	}
}