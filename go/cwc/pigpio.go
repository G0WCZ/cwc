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
	"github.com/golang/glog"
	"github.com/stianeikeland/go-rpio"
	"log"
	"strconv"
)

// PWM settings
const OnDutyCycle = uint32(1)
const PWMCycleLength = uint32(32)

type PiGPIO struct {
	config ConfigMap
	output rpio.Pin
	input rpio.Pin
	pwm bool
	pwmOut rpio.Pin
}

func NewPiGPIO() *PiGPIO {
	pigpio := PiGPIO{
		config: make(ConfigMap),
		pwm: false,
	}
	return &pigpio
}

func (g *PiGPIO) Open() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	sFreq, err := strconv.Atoi(g.config[Sidetonefreq])

	if err != nil {
		log.Fatalf("Bad sidetone frequency")
	}

	glog.Infof("setting sidetone to %d", sFreq)

	// PCM output
	if (sFreq > 0) {
		pcmPinNo, err := strconv.Atoi(g.config[Pcmout])
		if err != nil {
			log.Fatalf("bad pcmout value: %s", g.config[Pcmout])
		}
		g.pwm = true
		g.pwmOut = rpio.Pin(pcmPinNo)
		g.pwmOut.Mode(rpio.Pwm)
		g.pwmOut.Freq(sFreq * 32)
		g.pwmOut.DutyCycle(0, 32)
	}

	outPin, err := strconv.Atoi(g.config[Keyout])
	if err != nil {
		log.Fatalf("bad key value: %s", g.config[Keyout])
	}

	inPin, err := strconv.Atoi(g.config[Keyin])
	if err != nil {
		log.Fatalf("bad keyin value %s", g.config[Keyin])
	}

	// Pin output
	g.output = rpio.Pin(outPin)
	g.output.Output()
	g.output.Low()

    g.input = rpio.Pin(inPin) // header pin 13 BCM27
    g.input.Input()
    g.input.PullUp()

    return nil
}

func (g *PiGPIO) SetConfig(key string, value string) {
	g.config[key] = value
}

func (g *PiGPIO) ConfigMap() ConfigMap {
	return g.config
}

func (g *PiGPIO) Bit() bool {
	if g.input.Read() == rpio.High {
		return false
	} else {
		return true
	}
}

func (g *PiGPIO) SetBit(bit0 bool) {
	if bit0 {
		g.output.High()
		g.SetToneOut(true)
	} else {
		g.output.Low()
		g.SetToneOut(false)
	}
}

func (g *  PiGPIO) SetToneOut(v bool) {
	if g.pwm {
		var dutyLen uint32

		if v {
			dutyLen = OnDutyCycle
		} else {
			dutyLen = 0
		}
		g.pwmOut.DutyCycle(dutyLen, PWMCycleLength)
	}
}

func (g *PiGPIO) Close() {
	// pass
}

