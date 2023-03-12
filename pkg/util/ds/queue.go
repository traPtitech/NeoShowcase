package ds

import (
	"sync"

	"github.com/samber/lo"
)

type Queue[T any] struct {
	data  []T
	mutex sync.RWMutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Push(elt T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.data = append(q.data, elt)
}

func (q *Queue[T]) Pop() (elt T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.data) == 0 {
		return
	}
	elt = q.data[0]
	q.data = q.data[1:]
	return
}

func (q *Queue[T]) DeleteIf(predicate func(elt T) bool) (deleted bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	lenBefore := len(q.data)
	q.data = lo.Filter(q.data, func(elt T, i int) bool {
		return !predicate(elt)
	})
	return lenBefore != len(q.data)
}

func (q *Queue[T]) Len() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.data)
}
