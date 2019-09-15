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
package hw

import (
	"github.com/G0WCZ/cwc/config"
	"github.com/stianeikeland/go-rpio"
)

type GPIOOut struct {
	Config      *config.Config
	output      rpio.Pin
	adapterName string
}

func NewGPIOOut(config *config.Config, adapterName string) MorseOut {
	gpio_out := GPIOOut{
		Config:      config,
		adapterName: adapterName,
	}
	return &gpio_out
}

func (g *GPIOOut) Open() error {
	return g.initPorts(g.Config)
}

func (g *GPIOOut) ConfigChanged() error {
	return g.initPorts(g.Config)
}

func (g *GPIOOut) initPorts(config *config.Config) error {
	g.output = rpio.Pin(config.GPIOPins.KeyOut)
	g.output.Output()
	g.output.Low()

	return nil
}

func (g *GPIOOut) Close() error {
	return nil
}

func (g *GPIOOut) SetBit(bit bool) {
	if bit {
		g.output.High()
	} else {
		g.output.Low()
	}
}

func (g *GPIOOut) SetToneOut(bool) {
	panic("implement me")
}

func (g *GPIOOut) SetStatusLED(bool) {
	panic("implement me")
}
