package bitoip

import (
	"context"
	"github.com/golang/glog"
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

	if err != nil {
		glog.Fatalf("Can not open local connection: %v", err)
		return
	}
	defer conn.Close()

	glog.V(2).Infof("UDP Rx connection: %v", conn)

	buffer := make([]byte, maxBufferSize)
	doneChan := make(chan error, 1)

	go func() {
		for {
			n, addr, err := conn.ReadFromUDP(buffer)

			if err != nil {
				doneChan <- err
				return
			}

			glog.V(2).Infof("packet rx: %#v", buffer[0:n])

			verb, payload := DecodePacket(buffer)

			glog.V(2).Infof("udp rx got %v", payload)

			messages <- RxMSG{verb, payload, *addr}
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}
}
