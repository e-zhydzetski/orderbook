package threads

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestThreads(t *testing.T) {
	const parallel = 100
	const respDelay = 5

	log.Printf("OS: %s Arch: %s, GO_MAX_PROCS: %d, parallel: %d, respDelay: %d",
		runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(-1), parallel, respDelay)
	PrintStats()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go MonitorStats(ctx, 2*time.Second)

	var mx sync.Mutex
	responses := map[int]int{}

	var wg sync.WaitGroup
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("https://httpbin.org/delay/%d", respDelay))
			if err != nil {
				t.Error(err)
				return
			}
			mx.Lock()
			responses[resp.StatusCode]++
			mx.Unlock()
			_ = resp.Body.Close()
		}()
	}
	wg.Wait()

	PrintStats()

	log.Print("Responses:", responses)
}
