package eventbus

import (
	"context"

	"github.com/leandro-lugaresi/hub"
)

type local struct {
	bus *hub.Hub
}

type localSubscription struct {
	local *local
	sub   hub.Subscription
	c     chan *Event
	close chan struct{}
}

func newLocalSubscription(local *local, sub hub.Subscription) *localSubscription {
	s := &localSubscription{
		local: local,
		sub:   sub,
		c:     make(chan *Event),
		close: make(chan struct{}),
	}
	go s.relay()
	return s
}

func (s *localSubscription) relay() {
	for {
		select {
		case <-s.close:
			s.local.bus.Unsubscribe(s.sub)
			close(s.c)
			return
		case ev := <-s.sub.Receiver:
			s.c <- &Event{Type: ev.Name, Body: Fields(ev.Fields)}
		}
	}
}

func (s *localSubscription) Chan() <-chan *Event {
	return s.c
}

func (s *localSubscription) Unsubscribe() {
	close(s.close)
}

func NewLocal(hub *hub.Hub) Bus {
	return &local{bus: hub}
}

func (l *local) Publish(eventType string, body Fields) {
	l.bus.Publish(hub.Message{
		Name:   eventType,
		Fields: hub.Fields(body),
	})
}

func (l *local) Subscribe(events ...string) Subscription {
	sub := l.bus.Subscribe(100, events...)
	return newLocalSubscription(l, sub)
}

func (l *local) Close(ctx context.Context) error {
	l.bus.Close()
	return nil
}
