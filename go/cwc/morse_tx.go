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
