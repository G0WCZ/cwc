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
	"../bitoip"
	"gotest.tools/assert"
	"testing"
	"time"
)


func TestBuildPayload(t *testing.T) {
	var events = make([]Event, 0, MaxEvents)
	var now = time.Now()
	events = append(events,
		Event {
			now,
			bitoip.BitEvent(bitoip.BitOn),
		},
		)
	cep := BuildPayload(events)
	assert.Equal(t, cep.StartTimeStamp, now.UnixNano())
	assert.Equal(t, cep.BitEvents[0].BitEvent, bitoip.BitEvent(bitoip.BitOn|bitoip.LastEvent))
	assert.Equal(t, cep.BitEvents[0].TimeOffset, uint32(0))
}

func TestFlush0Events(t *testing.T) {
	c := make(chan bitoip.CarrierEventPayload)
	events := Flush([]Event{}, c)
	assert.Equal(t, len(events), 0)
}

func TestFlushMultipleEvents(t *testing.T) {
	c := make(chan bitoip.CarrierEventPayload)
	t1 := time.Now()
	t2 := time.Now()
	events := []Event{
		Event{t1, bitoip.BitOn},
		Event{t2, bitoip.BitOff},
	}
	go func() {
		e := Flush(events, c)
		assert.Equal(t, len(e), 0)
	}()

	cbe := <- c
	assert.Equal(t, cbe.Channel, uint16(0))
	assert.Equal(t, cbe.CarrierKey, uint16(0))
	assert.Equal(t, cbe.StartTimeStamp, t1.UnixNano())
	assert.Equal(t, len(cbe.BitEvents), 35)
}

