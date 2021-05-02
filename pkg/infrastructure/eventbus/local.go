package eventbus

import "github.com/leandro-lugaresi/hub"

type local struct {
	bus *hub.Hub
}

func (l *local) Publish(eventType string, body Fields) {
	l.bus.Publish(hub.Message{
		Name:   eventType,
		Fields: hub.Fields(body),
	})
}

func NewLocal(hub *hub.Hub) Bus {
	return &local{bus: hub}
}
