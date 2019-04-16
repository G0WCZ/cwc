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
	Close()
}

type ConfigMap map[string]string