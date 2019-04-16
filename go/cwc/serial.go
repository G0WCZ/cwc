package cwc

import (
	"github.com/stianeikeland/go-rpio"
	"go.bug.st/serial.v1"
	"log"
	"strings"
)

type SerialIO struct {
	config ConfigMap
	port serial.Port
}

func (s *SerialIO) Open() error {
	serialDevice := s.config["serialDevice"]

	log.Printf("Opening serial port %s", serialDevice)

	mode := &serial.Mode{}

	port, err := serial.Open(serialDevice, mode)
	s.port = port

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *SerialIO) SetConfig(key string, value string) {
	s.config[key] = value
}

func (s *SerialIO) ConfigMap() ConfigMap {
	return s.config
}

func (s *SerialIO) Bit() uint8 {
	info := s.port.Info
	info.
	if s.port.Description()

	}
	if s.input.Read() ==  rpio.High {
		return 0x01
	} else {
		return 0x00
	}
}

func (s *SerialIO) SetBit(bit0 uint8) {
	if strings.EqualFold("RTS", s.config["keyIn"]) {
		s.port.SetRTS(bits)
	} sles
	if bit0 & 0x01 > 0 {
		s.output.High()
	} else {
		s.output.Low()
	}
}


func (s *SerialIO) Close() {
	// pass
}


