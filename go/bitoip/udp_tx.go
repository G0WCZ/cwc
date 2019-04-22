package bitoip

import (
	"log"
	"net"
)


func UDPTx(verb MessageVerb, payload Payload, resolvedAddress *net.UDPAddr) {

	messagePayload := EncodePayload(verb, payload)
	connection := UDPConnection()
	log.Printf("udp connection %v", connection)
	n, err := connection.WriteToUDP(messagePayload, resolvedAddress)
	log.Printf("sent: %d, err %v", n, err)
}

