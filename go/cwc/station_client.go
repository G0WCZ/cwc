package cwc

import (
	"context"
	"github.com/golang/glog"
	"net"
	"../bitoip"
	"time"
)

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events

func StationClient(ctx context.Context, cqMode bool, addr string, morseIO IO, testFeedback bool, echo bool) {

	resolvedAddress, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		glog.Errorf("Error resolving address %s %v", addr, err)
		return
	}
	toSend := make(chan bitoip.CarrierEventPayload)
	toMorse := make(chan bitoip.RxMSG)

	// Morse receiver
	go RunMorseRx(ctx, morseIO, toSend, echo)

	localRxAddress, err := net.ResolveUDPAddr("udp", "0.0.0.0:8873")

	if err != nil {
		glog.Errorf("Can't allocate local address: %v", err)
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
			glog.V(2).Infof("carrier event payload to send: %v", cep)
			// TODO fill in some channel details
			bitoip.UDPTx(bitoip.CarrierEvent, cep, resolvedAddress)
			if testFeedback {
				QueueForTransmit(&cep)
			}

		case tm := <- toMorse:
			switch tm.Verb {
			case bitoip.CarrierEvent:
				glog.V(2).Infof("carrier events to morse: %v", tm)
				QueueForTransmit(tm.Payload.(*bitoip.CarrierEventPayload))

			case bitoip.ListenConfirm:
				glog.V(2).Infof("listen confirm: %v", tm)
				lc := tm.Payload.(*bitoip.ListenConfirmPayload)
				SetCarrierKey(lc.CarrierKey)
			}
		}
	}
}

