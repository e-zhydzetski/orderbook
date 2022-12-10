package threads

import (
	"context"
	"log"
	"runtime/pprof"
	"time"
)

var threadProfile = pprof.Lookup("threadcreate")
var goroutineProfile = pprof.Lookup("goroutine")

func PrintStats() {
	log.Printf("threads: %d, goroutines: %d, fd: %d",
		threadProfile.Count(), goroutineProfile.Count(), getFileDescriptorsCount())
}

func MonitorStats(ctx context.Context, interval time.Duration) {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			PrintStats()
		}
	}
}
