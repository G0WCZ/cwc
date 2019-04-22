package cwc

import (
	"context"
	"log"
	"net"
	"../bitoip"
	"time"
)

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events

func StationClient(ctx context.Context, cqMode bool, addr string, morseIO IO, testFeedback bool) {

	resolvedAddress, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		log.Printf("Error resolving address %s %v", addr, err)
		return
	}
	toSend := make(chan bitoip.CarrierEventPayload)
	toMorse := make(chan bitoip.RxMSG)

	// Morse receiver
	go RunMorseRx(ctx, morseIO, toSend)

	localRxAddress, err := net.ResolveUDPAddr("udp", "0.0.0.0:8873")

	if err != nil {
		log.Printf("Can't allocate local address: %v", err)
	}

	// UDP Receiver
	go bitoip.UDPRx(ctx, localRxAddress, toMorse)

	if !cqMode {
		time.Sleep(time.Second * 5)
		// TODO: full reflector mode implementation
		// Reflector mode setup
		// 1/ time sync with server
		// 2/ set callsign
		// 3/ list channels
		// 4/ suscribe channel(s)
		// 5/ save carrier id

		var csBase [16]byte
		bitoip.UDPTx(bitoip.ListenRequest, bitoip.ListenRequestPayload{
			0,
			 csBase,
			},
			resolvedAddress,
		)
	}
	for {
		select {
		case <- ctx.Done():
			return

		case cep := <- toSend:
			log.Printf("carrier event payload to send: %v", cep)
			// TODO fill in some channel details
			bitoip.UDPTx(bitoip.CarrierEvent, cep, resolvedAddress)
			if testFeedback {
				QueueForTransmit(&cep)
			}

		case tm := <- toMorse:
			switch tm.Verb {
			case bitoip.CarrierEvent:
				log.Printf("carrier events to morse: %v", tm)
				QueueForTransmit(tm.Payload.(*bitoip.CarrierEventPayload))
			}
		}
	}
}

