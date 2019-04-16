package cwc

/**
 * Specifies a device interface interface that can have
 * specific implementations for different I/O setups eg. gpio pins, serial port, audio out generator
 */

type IO interface {
	Open() error
	SetConfig(string, string)
	ConfigMap() ConfigMap
	Bit() uint8
	SetBit(bit0 uint8)
	Close()
}

type ConfigMap map[string]string