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

// Null I/O type
// useful for testing, doesn't aim to do anything
// with actual hardware

type NullIO struct {
	config ConfigMap
}

const BitIn = "bitin"
const BitOut = "bitout"
const ToneOut = "toneout"

func NewNullIO() *NullIO {
	return &NullIO{
		config: make(ConfigMap),
	}
}

func (n *NullIO) Open() error {
	return nil
}

func (n *NullIO) SetConfig(key string, value string) {
	n.config[key] = value
}

func (n *NullIO) ConfigMap() ConfigMap {
	return n.config
}

func (n *NullIO) Bit() bool {
	return n.config[BitIn] == "true"
}

func (n *NullIO) SetBit(b bool) {
	if b {
		n.config[BitOut] = "true"
	} else {
		n.config[BitOut] = "false"
	}
}

func (n * NullIO) SetToneOut(b bool) {
	if b {
		n.config[ToneOut] = "true"
	} else {
		n.config[ToneOut] = "false"
	}
}

func (* NullIO) Close() {}


