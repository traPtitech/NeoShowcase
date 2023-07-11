package domain

import (
	"context"
	"time"
)

type ContainerLog struct {
	Time time.Time
	Log  string
}

type ContainerLogger interface {
	Get(ctx context.Context, app *Application, before time.Time) ([]*ContainerLog, error)
	Stream(ctx context.Context, app *Application) (<-chan *ContainerLog, error)
}
