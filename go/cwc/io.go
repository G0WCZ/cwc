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