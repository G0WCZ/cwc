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
package cwc

import (
	"go.bug.st/serial.v1"
	"github.com/golang/glog"
	"strings"
)

type SerialIO struct {
	config ConfigMap
	port serial.Port
	useRTS bool
	useCTS bool
}

func NewSerialIO() *SerialIO {
	serialIO := SerialIO{
		config: make(ConfigMap),
		port: nil,
		useRTS: true,
		useCTS: true,
	}
	return &serialIO
}

func (s *SerialIO) Open() error {
	serialDevice := s.config["serialDevice"]

	glog.Infof("Opening serial port %s", serialDevice)

	mode := &serial.Mode{}

	port, err := serial.Open(serialDevice, mode)
	s.port = port

	if err != nil {
		glog.Fatalf("Can not open serial port: %v", err)
	}

	s.useRTS = strings.EqualFold(s.config[Keyout], "RTS")
	s.useCTS = strings.EqualFold(s.config[Keyin], "CTS")

	return nil
}

func (s *SerialIO) SetConfig(key string, value string) {
	s.config[key] = value
}

func (s *SerialIO) ConfigMap() ConfigMap {
	return s.config
}

func (s *SerialIO) Bit() bool {
	bits, err := s.port.GetModemStatusBits()
	if err != nil {
		glog.Fatalf("Port bit read failed %v", err)
		return false
	}
	if (s.useCTS) {
		return bits.CTS
	} else {
		return bits.DSR
	}
}

func (s *SerialIO) SetBit(bit bool) {
	var err error
	if s.useRTS {
		err = s.port.SetRTS(bit)
	} else {
		err = s.port.SetDTR(bit)
	}
	if err != nil {
		glog.Fatalf("port bit set failed: %v", err)
	}
}

func (s *SerialIO) SetToneOut(_ bool) {}

func (s *SerialIO) Close() {
	s.port.Close()
}


