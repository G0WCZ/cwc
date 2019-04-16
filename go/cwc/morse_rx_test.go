package cwc

import (
	"testing"
		"gotest.tools/assert"
	"time"
	"../bitoip"
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
	assert.Equal(t, cep.StartTimeStamp, now.Unix())
	assert.Equal(t, cep.BitEvents[0].BitEvent, bitoip.BitEvent(bitoip.BitOn))
	assert.Equal(t, cep.BitEvents[0].TimeOffset, uint32(0))
}
