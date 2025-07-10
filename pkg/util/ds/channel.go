package ds

import (
	"sync"
	"time"

	"github.com/bep/debounce"
)

func DebouncedChan[T any](ch <-chan T, period time.Duration) <-chan T {
	debouncedCh := make(chan T)

	var next T
	var lock sync.Mutex
	debounced := debounce.New(period)
	send := func() {
		lock.Lock()
		debouncedCh <- next
		lock.Unlock()
	}

	go func() {
		for data := range ch {
			lock.Lock()
			next = data
			lock.Unlock()
			debounced(send)
		}
	}()
	return debouncedCh
}
