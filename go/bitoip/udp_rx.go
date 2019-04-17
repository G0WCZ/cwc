package bitoip

import (
	"log"
	"fmt"
	"context"
	"net"
)

const maxBufferSize = 508

type RxMSG struct {
	Verb MessageVerb
	Payload Payload
	SrcAddress net.Addr
}

func UDPRx(ctx context.Context, address net.UDPAddr, messages chan RxMSG) {
	pc, err := net.ListenPacket("udp", address.String())

	if err != nil {
		return
	}

	defer pc.Close()

	buffer := make([]byte, maxBufferSize)
	doneChan := make(chan error, 1)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)

			if err != nil {
				doneChan <- err
				return
			}

			log.Printf("packet rx: %#v", buffer[0:n])

			verb, payload := DecodePacket(buffer)

			messages <- RxMSG{verb, payload, addr}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

}
