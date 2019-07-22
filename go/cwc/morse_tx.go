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

import "context"

import "../bitoip"

// This is the morse sender
// This sends output bits/streams

func MorseTx(ctx context.Context, morseToSend chan bitoip.CarrierEventPayload, config *Config) {

}

// Resolve config into actual outputs that are enabled
// For example, might be a decoder and a serial bit output
func ConfigureMorseOutputs(config *Config) {

}

func OpenOutputs() {

}

func CloseOutputs() {

}
