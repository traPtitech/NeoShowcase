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
	LogLimit() int
	Get(ctx context.Context, app *Application, before time.Time, limit int) ([]*ContainerLog, error)
	Stream(ctx context.Context, app *Application, begin time.Time) (<-chan *ContainerLog, error)
}
