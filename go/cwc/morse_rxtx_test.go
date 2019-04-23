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
	assert.Equal(t, len(cbe.BitEvents), 37)
}

