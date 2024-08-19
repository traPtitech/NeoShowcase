package domain

import "time"

type RuntimeImage struct {
	BuildID   string
	Size      int64
	CreatedAt time.Time
}

func NewRuntimeImage(buildID string, size int64) *RuntimeImage {
	return &RuntimeImage{
		BuildID:   buildID,
		Size:      size,
		CreatedAt: time.Now(),
	}
}
