package bitoip

import (
	"github.com/golang/glog"
	"net"
)


func UDPTx(verb MessageVerb, payload Payload, resolvedAddress *net.UDPAddr) {

	messagePayload := EncodePayload(verb, payload)
	connection := UDPConnection()
	n, err := connection.WriteToUDP(messagePayload, resolvedAddress)
	glog.V(2).Infof("sent udp: %d, err %v", n, err)
}

