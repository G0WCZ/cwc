package bitoip

import (
	"log"
	"net"
)


func UDPTx(verb MessageVerb, payload Payload, address string, local_address *net.UDPAddr) {
	resolved_address, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		log.Printf("Error resolving address %s %v", address, err)
		return
	}

	connection, err := net.DialUDP("udp", local_address, resolved_address)

	if err != nil {
		log.Printf("UDP Dial Error %v", err)
		return
	}

	messagePayload := EncodePayload(verb, payload)

	connection.Write(messagePayload)
}

