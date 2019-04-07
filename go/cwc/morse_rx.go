package cwc

import (
	"log"
	"time"
)
import "../bitoip"

/**
 * Morse receiver
 *
 * Takes incoming morse (a bit going high and low) and turns it into
 * CarrierBitEvents to send.
 *
 * Based on a regular tick that samples the input and builds a buffer
 */

const Ms = int64(1e6)
const Us = int64(1000)
const DefaultTickTime = time.Duration(1 * Ms)
const DefaultSendWait = time.Duration(100 * Ms)
const MaxEvents = 100

var TickTime = DefaultTickTime
var SendWait = DefaultSendWait

var LastBit byte = 0x00

type Event struct {
	startTime time.Time
	bitEvent bitoip.BitEvent
}

var events = make([]Event, 0, MaxEvents)

var ticker *time.Ticker
var done = make(chan bool)

func SetTickTime(tt time.Duration) {
	TickTime = tt
}

func SetSendWait(sw time.Duration) {
	SendWait = sw
}

var morseIO = PiGPIO{}

func RunRx() {
	LastBit = 0 // make sure turned off to begin -- the default state
	ticker = time.NewTicker(TickTime)
	Startup()

	for {
		select {
		case <- done:
			ticker.Stop()
			return

		case t := <-ticker.C:
			Sample(t)
		}
	}

}

func Stop() {
	done <- true
	LastBit = 0
	morseIO.Close()
}

func Startup() {
	err := morseIO.Open()
	if err != nil {
		log.Fatalf("Can't access GPIO: %s", err)
	}
}

// Sample input a
func Sample(t time.Time) {
	bit := morseIO.Bit()
	if bit != LastBit {
		// change so record it
		LastBit = bit
		events = append(events, Event{t, bitoip.BitEvent(bit) })
		if  len(events) >= MaxEvents {
			Flush()
			return
		}
	}
	if len(events)> 0 && t.Sub(events[0].startTime) >= DefaultSendWait {
		Flush()
	}
}


// Flush events into an output stream
func Flush() {
	log.Printf("Sending %v", events)
	events = events[:0]
}