package domain

import (
	"context"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

type Fields map[string]interface{}

type Bus interface {
	Publish(eventType string, body Fields)
	Subscribe(events ...string) Subscription
	Close(ctx context.Context) error
}

type Subscription interface {
	Chan() <-chan *Event
	Unsubscribe()
}

type Event struct {
	Type string
	Body Fields
}
