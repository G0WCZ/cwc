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

	go bitoip.UDPRx(ctx, serverAddress, messages)

	for {
		select {
		case <- ctx.Done():
			return
		case m := <- messages:
			Handler(serverAddress, m)
		}
	}
}

