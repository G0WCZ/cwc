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

/**
 * Specifies a device interface interface that can have
 * specific implementations for different I/O setups eg. gpio pins, serial port, audio out generator
 */

type IO interface {
	Open() error
	SetConfig(string, string)
	ConfigMap() ConfigMap
	Bit() bool
	SetBit(bool)
	SetToneOut(bool)
	Close()
}

type ConfigMap map[string]string

// Config consts
const Keyin = "keyin"
const Keyout = "keyout"
const Pcmout = "pcmout"
const Sidetonefreq = "sidetonefreq"