package ds

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMutex_Lock(t *testing.T) {
	m := NewMutex[string]()

	start := time.Now()
	const n = 10

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			m.Lock("k1")
			time.Sleep(100 * time.Millisecond)
			m.Unlock("k1")
		}()
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			m.Lock("k2")
			time.Sleep(100 * time.Millisecond)
			m.Unlock("k2")
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)
	assert.True(t, 1000 < elapsed.Milliseconds() && elapsed.Milliseconds() < 1200)
}

func TestMutex_TryLock(t *testing.T) {
	t.Run("already locked", func(t *testing.T) {
		m := NewMutex[string]()

		m.Lock("k1")
		m.Lock("k2")
		assert.False(t, m.TryLock("k1"))
		assert.False(t, m.TryLock("k2"))
	})

	t.Run("success", func(t *testing.T) {
		m := NewMutex[string]()

		assert.True(t, m.TryLock("k1"))
		assert.False(t, m.TryLock("k1"))
		assert.True(t, m.TryLock("k2"))
		assert.False(t, m.TryLock("k2"))
	})
}

func BenchmarkMutex(b *testing.B) {
	m := NewMutex[string]()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Lock("k1")
		m.Unlock("k1")
	}
}
