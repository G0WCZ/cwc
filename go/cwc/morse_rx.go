package cwc

import "context"

import "../bitoip"

// This is the morse receiver
// This polls morse input "hardware" for data

func MorseRx(ctx context.Context, morseReceived chan bitoip.CarrierEventPayload, config *Config) {

}

// Resolve config into actual inputs that are enabled
// For example, might be a key/keyer and text input
func ConfigureMorseInputs(config *Config) {

}

func OpenInputs() {

}

func CloseInputs() {

}

// Keyer support
