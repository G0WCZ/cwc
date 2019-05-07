package main

import (
	"context"
	"github.com/golang/glog"
	"time"
)

const StationTimeout = 5 * time.Minute

func Supervisor(ctx context.Context) {
	tick := time.Tick(60 * time.Second)

	for {
		select {
		case t := <- tick:
			SuperviseReflector(t)
		case <-ctx.Done():
			return
		}
	}
}

// Supervise Reflector -- generally tidy up
func SuperviseReflector(t time.Time) {
	removedCount := SuperviseChannels(t, StationTimeout)
	glog.Infof("Supervisor removed %d stale stations.", removedCount)
}