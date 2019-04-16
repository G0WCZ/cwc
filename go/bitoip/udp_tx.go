package bitoip

import (
	"net"
)


func UDPTx(verb MessageVerb, payload Payload, address string, local_address *net.UDPAddr) {
	resolved_address, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		return
	}

	connection, err := net.DialUDP("udp", local_address, resolved_address)

	if err != nil {
		return
	}

	messagePayload := EncodePayload(verb, payload)

	connection.Write(messagePayload)
}

