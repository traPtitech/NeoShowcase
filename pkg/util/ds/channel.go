package ds

import (
	"context"
	"sync"
	"time"

	"github.com/motoki317/go-stabilize"
)

func StabilizeChan[T any](ch <-chan T, period time.Duration) <-chan T {
	stabilizedCh := make(chan T)

	var next T
	var lock sync.Mutex
	stabilizer := stabilize.NewStabilizer(context.Background(), period, func() {
		lock.Lock()
		stabilizedCh <- next
		lock.Unlock()
	})

	go func() {
		for data := range ch {
			lock.Lock()
			next = data
			lock.Unlock()
			stabilizer.Trigger()
		}
	}()
	return stabilizedCh
}
