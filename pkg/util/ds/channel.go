package ds

import (
	"sync"
	"time"

	"github.com/boz/go-throttle"
)

func ThrottleChan[T any](ch <-chan T, period time.Duration) <-chan T {
	throttledCh := make(chan T)

	var next T
	var lock sync.Mutex
	throttledSend := throttle.ThrottleFunc(period, true, func() {
		lock.Lock()
		throttledCh <- next
		lock.Unlock()
	})

	go func() {
		for data := range ch {
			lock.Lock()
			next = data
			lock.Unlock()
			throttledSend.Trigger()
		}
	}()
	return throttledCh
}
