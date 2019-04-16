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

	bitoip.UDPRx(ctx, *serverAddress, Handler)
}

/**
	Handle an incoming message to the reflector
 */
func Handler(verb bitoip.MessageVerb, payload bitoip.Payload, remote_address net.Addr) {
	switch verb {
	case bitoip.EnumerateChannels:
		responsePayload := new(bitoip.ListChannelsPayload)
		copy(responsePayload.Channels[:], ChannelIds())
		bitoip.UDPTx(bitoip.ListChannels,
					 responsePayload,
					 remote_address.String(),
					 serverAddress)

	case bitoip.CarrierEvent:

	}
}