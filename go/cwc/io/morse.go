package io

import "time"

import "../../bitoip"

// Absolute time morse event
type TimedBitEvent struct {
	startTime time.Time
	bitEvent  bitoip.BitEvent
}

// slice of events
type MorseEvents []TimedBitEvent
