package main

import (
	"context"
	"time"
)

func Supervisor(ctx context.Context) {
	tick := time.Tick(time.Minute)
	for {
		select {
		case t := <- tick
			SuperviseReflector(t)
		case <-ctx.Done():
			return
		}
	}
}

func SuperviseReflectot(t time.Time) {

}