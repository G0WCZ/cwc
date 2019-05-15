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

