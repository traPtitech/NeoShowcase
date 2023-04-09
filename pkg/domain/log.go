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
	Get(ctx context.Context, appID string, before time.Time) ([]*ContainerLog, error)
	Stream(ctx context.Context, appID string, after time.Time) (<-chan *ContainerLog, error)
}
