package domain

import (
	"sync"

	"github.com/samber/lo"
)

// PubSub is a simple pub-sub for single event type.
// A zero-value PubSub is ready for use.
type PubSub[E any] struct {
	subs []chan<- E
	lock sync.Mutex
}

func (p *PubSub[E]) Publish(event E) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, sub := range p.subs {
		select {
		case sub <- event:
		default:
		}
	}
}

func (p *PubSub[E]) Subscribe() (sub <-chan E, unsub func()) {
	p.lock.Lock()
	defer p.lock.Unlock()

	ch := make(chan E)
	p.subs = append(p.subs, ch)
	return ch, func() {
		p.lock.Lock()
		p.subs = lo.Without(p.subs, ch)
		p.lock.Unlock()
	}
}
