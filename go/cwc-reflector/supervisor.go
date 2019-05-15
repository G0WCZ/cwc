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