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
package bitoip

import (
	"context"
	"github.com/G0WCZ/cwc/cwcpb"
	"github.com/golang/glog"
	"net"
	"time"
)

const maxBufferSize = 508

var conn *net.UDPConn

type (
	RxMSG struct {
		Message    *cwcpb.CWCMessage
		SrcAddress net.UDPAddr
		RxTime     int64
	}
)

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
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		glog.Fatalf("Can not zero read deadline: %v", err)
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

			now := time.Now().UnixNano()

			glog.V(2).Infof("packet rx: %#v", buffer[0:n])

			message := DecodeBuffer(buffer, n)

			glog.V(2).Infof("udp rx got %v", message)

			messages <- RxMSG{message, *addr, now}
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}
}
