package threads

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestThreads(t *testing.T) {
	requestString := os.Getenv("REQUEST_STR")
	if requestString == "" {
		t.Fatal("REQUEST_STR not defined")
	}

	const parallel = 1000

	log.Printf("OS: %s Arch: %s, GO_MAX_PROCS: %d, Requests: %d x %s",
		runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(-1), parallel, requestString)
	PrintStats()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go MonitorStats(ctx, 1*time.Second)

	var mx sync.Mutex
	responses := map[int]int{}

	var wg sync.WaitGroup
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Get(requestString)
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

	log.Printf("Response codes: %v", responses)
}
