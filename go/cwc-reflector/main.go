/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


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
package main

import (
	"../bitoip"
	"context"
	"flag"
	"github.com/golang/glog"
	"net"
	"os"
)

var serverAddress *net.UDPAddr

func main() {
	address := flag.String("address", "localhost:7388", "-address=host:port")

	flag.Parse()

	ReflectorServer(context.TODO(), *address)
}

func ReflectorServer(ctx context.Context, address string) {

	glog.Info(DisplayVersion())

	serverAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		glog.Fatalf("Can't use address %s: %s", address, err)
		os.Exit(1)
	}

	glog.Infof("Starting reflector on %s", address)

	messages := make(chan bitoip.RxMSG)

	go bitoip.UDPRx(ctx, serverAddress, messages)

	go APIServer(ctx, &channels, address)

	go Supervisor(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-messages:
			Handler(serverAddress, m)
		}
	}
}
