package cwc

type NullIO struct {
	config ConfigMap
}

const BitIn = "bitin"
const BitOut = "bitout"
const ToneOut = "toneout"

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


