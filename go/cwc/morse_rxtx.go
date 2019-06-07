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
package cwc

import (
	"github.com/golang/glog"
	"context"
	"sort"
	"sync"
	"time"
)
import "../bitoip"

/**
 * Morse hardware receiver and sender
 *
 * Takes incoming morse (a bit going high and low) and turns it into
 * CarrierBitEvents to send.
 *
 * Based on a regular tick that samples the input and builds a buffer
 */

const Ms = int64(1e6)
const Us = int64(1000)
const DefaultTickTime = time.Duration(5 * Ms)
const KeyerTickTime = time.Duration(1 * Ms)
const MaxSendTimespan = time.Duration(1000 * Ms)
const BreakinTime = time.Duration(100 * Ms)
const MaxEvents = 100

const (
	CHECK 		int = 0
	PREDOT 		int = 1
	PREDASH 	int = 2
	SENDDOT 	int = 3
	SENDDASH 	int = 4
	DOTDELAY 	int = 5
	DASHDELAY 	int = 6
	DOTHELD 	int = 7
	DASHHELD 	int = 8
	LETTERSPACE int = 9
	EXITLOOP 	int = 10
)

var TickTime = DefaultTickTime
var SendWait = MaxSendTimespan

var key_state = CHECK
var dot_memory = 0
var dash_memory = 0
var kdelay = 0
var dot_delay = 0
var dash_delay = 0


var LastBit bool = false

type Event struct {
	startTime time.Time
	bitEvent bitoip.BitEvent
}

var events = make([]Event, 0, MaxEvents)
var RxMutex = sync.Mutex{}

var ticker *time.Ticker
var done = make(chan bool)


func SetTickTime(tt time.Duration) {
	TickTime = tt
}

func SetSendWait(sw time.Duration) {
	SendWait = sw
}

var localEcho bool

var channelId bitoip.ChannelIdType

func SetChannelId(cId bitoip.ChannelIdType) {
	channelId = cId
}

func ChannelId () bitoip.ChannelIdType {
	return channelId
}

var carrierKey bitoip.CarrierKeyType

func SetCarrierKey(ck bitoip.CarrierKeyType) {
	carrierKey = ck
}

func CarrierKey() bitoip.CarrierKeyType {
	return carrierKey
}

// Clock offset from the reflector
var timeOffset = int64(0)

func SetTimeOffset(t int64) {
	timeOffset = t
}

// round trip delay calculated by the transport
var roundTrip = int64(0)

func SetRoundTrip(t int64) {
	roundTrip = t
}

// RunMorseRx sets the morse hardware (key) receiver going.  This sets up a timer
// to sample the morse input and runs it.


func RunMorseRx(ctx context.Context, morseIO IO, toSend chan bitoip.CarrierEventPayload, echo bool,
	channel bitoip.ChannelIdType, keyer bool) {
	localEcho = echo
	channelId = channel
	LastBit = false // make sure turned off to begin -- the default state
	if (keyer)
		TickTime = KeyerTickTime
	ticker = time.NewTicker(TickTime)

	Startup(morseIO)

	for {
		select {
		case <- done:
			ticker.Stop()
			return

		case t := <-ticker.C:
			if (keyer) 
				SampleKeyer(t, toSend, morseIO)
			else
				Sample(t, toSend, morseIO)
		}
	}
}

// Overload to ensure backwards compatibility
func RunMorseRx(ctx context.Context, morseIO IO, toSend chan bitoip.CarrierEventPayload, echo bool,
	channel bitoip.ChannelIdType) {
	RunMorseRx(ctx, morseIO, toSend, echo, channel, false)	
}

func Stop(morseIO IO) {
	done <- true
	LastBit = false
	morseIO.Close()
}

func Startup(morseIO IO) {
	err := morseIO.Open()
	if err != nil {
		glog.Fatalf("Can't access Morse hardware: %s", err)
	}
}

// Sample the input pin for the morse key.
// This is called (currently) every 5ms to look for a change in input pin.
//
// TODO should have some sort of back-off if not used recently for power saving
func Sample(t time.Time, toSend chan bitoip.CarrierEventPayload, morseIO IO) {

	TransmitToHardware(t, morseIO)

	rxBit := morseIO.Bit()
	if rxBit != LastBit {
		// change so record it
		LastBit = rxBit
		morseIO.SetToneOut(rxBit)

		var bit uint8 = 0

		if rxBit {
			bit = 1
		}

		// Append the event to the list of on/off events to be sent
		RxMutex.Lock()
		events = append(events, Event{t, bitoip.BitEvent(bit) })
		RxMutex.Unlock()

		// If we have too many events, then send a packet now
		if  (len(events) >= MaxEvents - 1) && (events[len(events)-1].bitEvent & bitoip.BitOn == 0) {
			events = Flush(events, toSend)
			return
		}
	}

	// Check to see if we've got an event buffer that needs sending because it has been
	// hanging around too long.
	if len(events)> 0 &&
		(events[len(events)-1].bitEvent & bitoip.BitOn == 0) &&
		((t.Sub(events[0].startTime) >= MaxSendTimespan) ||
		(t.Sub(events[len(events) - 1].startTime) >= BreakinTime)) {
		events = Flush(events, toSend)
	}
}

func clear_memory() {
	dot_memory = 0
	dash_memory = 0
}

func SetKeyerOut(int state) {
	if(keyer_out != state) {
		keyer_out = state
		morseIO.SetToneOut(bit);
		RxMutex.Lock()
		events = append(events, Event{t, bitoip.BitEvent(state) })
		RxMutex.Unlock()
		// If we have too many events, then send a packet now
		if  (len(events) >= MaxEvents - 1) && (events[len(events)-1].bitEvent & bitoip.BitOn == 0) {
			events = Flush(events, toSend)
			return
		}
	}	
}

func SampleKeyer(t time.Time, toSend chan bitoip.CarrierEventPayload, morseIO IO) {
	TransmitToHardware(t, morseIO)

	if (key_state != EXITLOOP) {
		switch (key_state) {
		case CHECK:
			if (morseIO.Dot()) {
				key_state = PREDOT
			} else if (morseIO.Dash()) {
				key_state = PREDASH
			} else {
				key_state = EXITLOOP
			}
		case PREDOT:
			clear_memory()
			key_state = SENDDOT
		case PREDASH:
			clear_memory()
			key_state = SENDDASH
		case SENDDOT:
			SetKeyerOut(1)
			if (kdelay == dot_delay) {
				kdelay = 0
				SetKeyerOut(0)
				key_state = DOTDELAY
			} else kdelay++
			if (cw_keyer_mode == "A")
				if ((!morseIO.Dot())&&(!morseIO.Dash()))
					dash_memory = 0
				else if (morseIO.Dash())
					dash_memory = 1
		case SENDDASH:
			SetKeyerOut(1)
			if (kdelay == dash_delay) {
				kdelay = 0
				SetKeyerOut(0)
				key_state = DASHDELAY
			} else kdelay++
			if (cw_keyer_mode == "A")
			if ((!morseIO.Dot())&&(!morseIO.Dash()))
				dot_memory = 0
			else if (morseIO.Dot())
				dot_memory = 1
		case DOTDELAY:
			if (kdelay == dot_delay) {
				kdelay = 0
				if (dash_memory)
					key_state = PREDASH
				else 
					key_state = DOTHELD
			} else kdelay++
			if (morseIO.Dash())
				dash_memory = 1
		case DASHDELAY:
			if (kdelay == dot_delay) {
				kdelay = 0
				if (dot_memory)
					key_state = PREDOT
				else 
					key_state = DASHHELD
			} else kdelay++
			if (morseIO.Dot())
				dash_memory = 1
		case DOTHELD:
			if (morseIO.Dot())
				key_state = PREDOT
			else if (morseIO.Dash())
				key_state = PREDASH
			else if (cw_keyer_spacing) {
				clear_memory()
				key_state = LETTERSPACE
			} else
				key_state = EXITLOOP
		case DASHHELD:
			if (morseIO.Dash())
				key_state = PREDASH
			else if (morseIO.Dot())
				key_state = PREDOT
			else if (cw_keyer_spacing) {
				clear_memory()
				key_state = LETTERSPACE
			} else
				key_state = EXITLOOP
		case LETTERSPACE:
			if (kdelay == 2 * dot_delay) {
				kdelay = 0
				if (dot_memory)
					key_state = PREDOT
				else if (dash_memory)
					key_state = PREDASH
				else
					key_state = EXITLOOP
			} else kdelay++
			if (morseIO.Dot()) dot_memory = 1
			if (morseIO.Dash()) dash_memory = 1
		case EXITLOOP:
			key_state = CHECK
		default:
			key_state = EXITLOOP
	}
	// Check to see if we've got an event buffer that needs sending because it has been
	// hanging around too long.
	if len(events)> 0 &&
		(events[len(events)-1].bitEvent & bitoip.BitOn == 0) &&
		((t.Sub(events[0].startTime) >= MaxSendTimespan) ||
		(t.Sub(events[len(events) - 1].startTime) >= BreakinTime)) {
		events = Flush(events, toSend)
	}

}


// Flush events and place in the toSend channel to wake up the UDP sender to
// transmit the packet.
func Flush(events []Event, toSend chan bitoip.CarrierEventPayload) []Event {
	glog.V(2).Infof("Flushing events %v", events)
	RxMutex.Lock()
	if len(events) > 0 {
		toSend <- BuildPayload(events)
		events = events[:0]
	}
	RxMutex.Unlock()
	return events
}

// Build a payload (CarrierEventPayload) of on and off events. Called from Flush() to
// make a packet ready to send.
func BuildPayload(events []Event) bitoip.CarrierEventPayload {
	baseTime := events[0].startTime.UnixNano()
	cep := bitoip.CarrierEventPayload{
		channelId,
		carrierKey,
		baseTime + timeOffset,
		[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{},
		time.Now().UnixNano(),
	}
	for i, event := range events {
		bit := event.bitEvent

		// mark last event this message
		if i == (len(events) - 1) {
			bit = bit | bitoip.LastEvent
		}

		cep.BitEvents[i] = bitoip.CarrierBitEvent{
			uint32(event.startTime.UnixNano() - baseTime),
			bit,
		}
	}
	return cep
}

/**
 * Transmitting morse out a gpio pin
 */

// Lock for the Morse Transmit Queue
var TxMutex = sync.Mutex{}

// Morse Transmit Queue
var TxQueue = make([]Event, 100)

// Queue this stuff for sending to hardware -- LED or relay or PWM
// by adding to queue that will be sent out based on the tick timing
func QueueForTransmit(carrierEvents *bitoip.CarrierEventPayload) {
	if (localEcho || (carrierEvents.CarrierKey != carrierKey)) &&
		carrierEvents.Channel == channelId {
		// compose into events
		newEvents := make([]Event, 0)

		// remove the calculated server time offset
		start := time.Unix(0, carrierEvents.StartTimeStamp - timeOffset + roundTrip + int64(MaxSendTimespan))

		for _, ce := range carrierEvents.BitEvents {
			newEvents = append(newEvents, Event{
				start.Add(time.Duration(ce.TimeOffset)),
				ce.BitEvent,
			})
			if (ce.BitEvent & bitoip.LastEvent) > 0 {
				break
			}
		}

		// Lock and append new events
		TxMutex.Lock()
		TxQueue = append(TxQueue, newEvents...)
		// then sort the output by time (this is probably super slow)
		sort.Slice(TxQueue, func(i, j int) bool { return TxQueue[i].startTime.Before(TxQueue[j].startTime) })
		TxMutex.Unlock()
	} else {
		// don't re-sound our own stuff if echo isn't turned on
		glog.V(2).Infof("ignoring own carrier")
	}
	glog.V(2).Infof("TXQueue is now: %v", TxQueue)
}

// When woken up  (same timer as checking for an incoming bit change)
// check to see if an output state change is needed and do it.
func TransmitToHardware(t time.Time, morseIO IO) {
	now := time.Now()

	// Lock
	TxMutex.Lock()

	// Change output if needed
	if len(TxQueue) > 0 && TxQueue[0].startTime.Before(now) {
		be := TxQueue[0].bitEvent
		morseIO.SetBit(!((be & bitoip.BitOn) == 0))
		TxQueue = TxQueue[1:]
	}

	// Unlock
	TxMutex.Unlock()
}
