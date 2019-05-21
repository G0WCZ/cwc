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
package main

import (
	"../bitoip"
	"../cwc"
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
)

const maxBufferSize = 508

func main() {
	var cqMode = flag.Bool("cq", false, "-cq for local mode")
	var configFile = flag.String("config", "/boot/cwc-station.txt", "-config <filename>")
	var refAddPtr = flag.String("ref", "", "-ref <host>:<port>")
	var serialDevice = flag.String("serial", "", "-serial=<serial-device-name>")
	var echo = flag.Bool("echo", false, "-echo turns on remote echo of all sent morse")
	var channel = flag.Int("ch", -1, "-ch <n> to connect to the channel n")
	var callsign = flag.String("de", "", "-de <callsign>")
	var noIO = flag.Bool("noio",false, "-noio uses fake morse IO connections")

	// parse Command line
	flag.Parse()

	// read Config file and defaults
	config := cwc.ReadConfig(*configFile)

	// Network mode
	if *cqMode {
		config.NetworkMode = "local"
	}
	// Reflector address
	if len(*refAddPtr) > 0 {
		config.ReflectorAddress = *refAddPtr
	}

	if len(*serialDevice) > 0 {
		config.SerialDevice = *serialDevice
		config.HardwareType = "Serial"
	}

	if *echo {
		config.RemoteEcho = true
	}

	if *channel >= 0 {
		config.Channel = bitoip.ChannelIdType(*channel)
	}

	if len(*callsign) > 0 {
		config.Callsign = *callsign
	}

	if *noIO {
		config.HardwareType = "None"
	}

	// context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println(DisplayVersion())

	glog.Info(DisplayVersion())

	// Morse Hardware
	var morseIO cwc.IO

	if config.HardwareType == "Serial" {
		morseIO = cwc.NewSerialIO(config)
	} else if config.HardwareType == "None" {
		morseIO = cwc.NewNullIO(config)
	} else {
		morseIO =  cwc.NewPiGPIO(config)
	}

	cwc.StationClient(ctx, config, morseIO)
}
