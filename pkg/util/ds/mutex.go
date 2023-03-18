package ds

import (
	"fmt"
	"sync"
)

// Mutex is a simple replacement for golang.org/x/sync/singleflight.
// Includes ideas from https://cs.opensource.google/go/x/sync/+/master:singleflight/singleflight.go;drc=30421366ff761c80b137fb5084b32278ed41fab0.
type Mutex[K comparable] struct {
	calls map[K]*call
	mu    sync.Mutex
}

func NewMutex[K comparable]() *Mutex[K] {
	return &Mutex[K]{calls: make(map[K]*call)}
}

type call struct {
	sem chan struct{}
	cnt int
}

func (m *Mutex[K]) Lock(key K) {
	m.mu.Lock()
	c, ok := m.calls[key]
	if !ok {
		c = &call{
			sem: make(chan struct{}, 1),
		}
		m.calls[key] = c
	}
	c.cnt++
	m.mu.Unlock()

	c.sem <- struct{}{}
}

func (m *Mutex[K]) TryLock(key K) (ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.calls[key]
	if ok {
		return false
	}

	c = &call{
		sem: make(chan struct{}, 1),
	}
	m.calls[key] = c
	c.cnt++
	c.sem <- struct{}{}
	return true
}

func (m *Mutex[K]) Unlock(key K) {
	m.mu.Lock()
	c, ok := m.calls[key]
	if !ok {
		panic(fmt.Sprintf("mapmutex: Unlock() called before Lock() on key: %v", key))
	}
	<-c.sem
	c.cnt--
	if c.cnt == 0 {
		delete(m.calls, key)
	}
	m.mu.Unlock()
}
