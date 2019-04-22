package bitoip

import (
	"log"
	"fmt"
	"context"
	"net"
)

const maxBufferSize = 508

var conn *net.UDPConn

type RxMSG struct {
	Verb MessageVerb
	Payload Payload
	SrcAddress net.UDPAddr
}

func UDPConnection() *net.UDPConn {
	return conn
}

func UDPRx(ctx context.Context, address *net.UDPAddr, messages chan RxMSG) {
	var err error
	conn, err = net.ListenUDP("udp", address)
	defer conn.Close()

	if err != nil {
		log.Panicf("Can not open local connection: %v", err)
		return
	}

	log.Printf("UDP Rx connection: %v", conn)

	buffer := make([]byte, maxBufferSize)
	doneChan := make(chan error, 1)

	go func() {
		for {
			n, addr, err := conn.ReadFromUDP(buffer)

			if err != nil {
				doneChan <- err
				return
			}

			log.Printf("packet rx: %#v", buffer[0:n])

			verb, payload := DecodePacket(buffer)

			log.Printf("got %v", payload)

			messages <- RxMSG{verb, payload, *addr}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}
}
