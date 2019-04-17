package cwc

import (
	"context"
	"log"
)
import "../bitoip"

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events

func StationClient(ctx context.Context, cqMode bool, addr string, morseIO IO) {

	// if talking to reflector, do some setup
	if !cqMode {
		// TODO:
		// Reflector mode setup
		// 1/ time sync with server
		// 2/ set callsign
		// 3/ list channels
		// 4/ suscribe channel(s)
		// 5/ save carrier id
	}

	toSend := make(chan bitoip.CarrierEventPayload)

	go RunRx(ctx, morseIO, toSend)

	for {
		select {
		case <- ctx.Done():
			return
		case cep := <- toSend:
			log.Printf("carrier event payload to send: %v", cep)
			// TODO fill in some channel details
			bitoip.UDPTx(bitoip.CarrierEvent, cep, addr,nil)
		}
	}

	toMorse := make(chan bitoip.CarrierEventPayload)

	go bitoip.UDPRx(ctx, addr, toMorse)

	for {
		select {
		case <-ctx.Done():
			return

		case tm := <- toMorse:
			log.Printf("carrier events to morse: %v", toMorse)

		}
	}
}

