package bitoip

import (
	"github.com/golang/glog"
	"net"
)


func UDPTx(verb MessageVerb, payload Payload, resolvedAddress *net.UDPAddr) {

	messagePayload := EncodePayload(verb, payload)
	connection := UDPConnection()
	glog.V(1).Infof("udp connection %v", connection)
	n, err := connection.WriteToUDP(messagePayload, resolvedAddress)
	glog.V(1).Infof("sent: %d, err %v", n, err)
}

