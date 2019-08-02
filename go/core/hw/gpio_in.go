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
	"strconv"
)

type GPIOIn struct {
	Config *config.Config
	input  rpio.Pin
}

func (g *GPIOIn) Open(config *config.Config) error {
	g.Config = config
	return g.initPorts(config)
}

func (g *GPIOIn) ConfigChanged() error {
	return g.initPorts(g.Config)
}

func (g *GPIOIn) initPorts(config *config.Config) error {
	inPin := strconv.Atoi(config.GPIOPins.)
	g.input = rpio.Pin(inPin)
	g.input.Input()
	g.input.PullUp()
	return nil
}

func (g *GPIOIn) Close() error {
	return nil
}
