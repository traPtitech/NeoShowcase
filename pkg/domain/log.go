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
	Stream(ctx context.Context, appID string) (<-chan *ContainerLog, error)
}
