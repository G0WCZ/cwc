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
package config

import (
	"fmt"
	"math/rand"

	"github.com/BurntSushi/toml"
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/golang/glog"
)

type Config struct {
	NetworkMode       string
	ReflectorAddress  string
	LocalPort         int
	MorseInHardware   []string
	MorseOutHardware  []string
	GeneralHardware   []string
	SidetoneEnable    bool
	SidetoneFrequency int
	RemoteEcho        bool
	Channel           bitoip.ChannelIdType
	Callsign          string
	GPIOPins          GPIOPins
	Serial            Serial
	Keyer             Keyer
	Encoder           Encoder
	Decoder           Decoder
}

const HWKeyTip = 17
const HWKeyRing = 27
const HWLEDStatus = 22
const HWLEDSignal = 23

type GPIOPins struct {
	KeyLeft   int
	KeyRight  int
	PWMA      int
	PWMB      int
	KeyOut    int
	StatusLED int
	SignalLED int
}

type Serial struct {
	Device   string
	KeyLeft  string
	KeyRight string
	KeyOut   string
}

type Keyer struct {
	Type    string
	Speed   int
	Weight  int
	Mode    int
	Reverse bool
}

type Encoder struct {
	Speed int //wpm
}

type Decoder struct {
	StartingSpeed int //wpm
}

var defaultConfig = Config{
	NetworkMode:       "Reflector",
	ReflectorAddress:  "cwc0.nodestone.io:7388",
	LocalPort:         5990,
	MorseInHardware:   []string{"keyer:pigpio"}, // GPIO or Serial or None
	MorseOutHardware:  []string{"pigpio"},
	SidetoneEnable:    true,
	SidetoneFrequency: 500,
	RemoteEcho:        false,
	Channel:           0,
	Callsign:          fmt.Sprintf("CWC%d", rand.Int31()),

	GPIOPins: GPIOPins{
		StatusLED: HWLEDStatus,
		SignalLED: HWLEDSignal,
		KeyLeft:   HWKeyTip,
		KeyRight:  HWKeyRing,
		PWMA:      13,
		PWMB:      12,
	},
	Serial: Serial{
		Device:   "/dev/ttysomething",
		KeyLeft:  "CTS",
		KeyRight: "",
		KeyOut:   "RTS",
	},

	Keyer: Keyer{
		Type:    "keyer",
		Speed:   20,
		Weight:  55,
		Mode:    1,
		Reverse: false,
	},
	Encoder: Encoder{
		Speed: 12,
	},
	Decoder: Decoder{
		StartingSpeed: 12,
	},
}

func ReadConfig(filename string) *Config {
	cfg := defaultConfig

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		glog.Warningf("Config file not found %s", filename)
	}

	return &cfg
}
