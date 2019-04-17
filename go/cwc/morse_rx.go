package cwc

import (
	"context"
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

var LastBit bool = false

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

func RunRx(ctx context.Context, morseIO IO, toSend chan bitoip.CarrierEventPayload) {
	LastBit = false // make sure turned off to begin -- the default state
	ticker = time.NewTicker(TickTime)

	Startup(morseIO)

	for {
		select {
		case <- done:
			ticker.Stop()
			return

		case t := <-ticker.C:
			Sample(t, toSend, morseIO)
		}
	}
}

func Stop(morseIO IO) {
	done <- true
	LastBit = false
	morseIO.Close()
}

func Startup(morseIO IO) {
	err := morseIO.Open()
	if err != nil {
		log.Fatalf("Can't access Morse hardware: %s", err)
	}
}

// Sample input
// TODO should have some sort of back-off if not used recently for power saving
func Sample(t time.Time, toSend chan bitoip.CarrierEventPayload, morseIO IO) {
	rxBit := morseIO.Bit()
	if rxBit != LastBit {
		// change so record it
		LastBit = rxBit

		var bit uint8 = 0

		if rxBit {
			bit = 1
		}
		events = append(events, Event{t, bitoip.BitEvent(bit) })
		if  len(events) >= MaxEvents {
			Flush(toSend)
			return
		}
	}
	if len(events)> 0 && t.Sub(events[0].startTime) >= DefaultSendWait {
		Flush(toSend)
	}
}


// Flush events into an output stream
func Flush(toSend chan bitoip.CarrierEventPayload) {
	log.Printf("Flushing events %v", events)
	if len(events) > 0 {
		toSend <- BuildPayload(events)
		events = events[:0]
	}
}

func BuildPayload(events []Event) bitoip.CarrierEventPayload {
	baseTime := events[0].startTime.Unix()
	cep := bitoip.CarrierEventPayload{
		0,
		0,
		baseTime,
		[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{},
	}
	for i, event := range events {
		cep.BitEvents[i] = bitoip.CarrierBitEvent{
			uint32(event.startTime.Unix() - baseTime),
			event.bitEvent,
		}
	}
	return cep
}
