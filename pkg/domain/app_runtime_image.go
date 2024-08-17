package domain

import "time"

type RuntimeImage struct {
	ID        string
	BuildID   string
	Size      int64
	CreatedAt time.Time
}

func NewRuntimeImage(buildID string, size int64) *RuntimeImage {
	return &RuntimeImage{
		ID:        NewID(),
		BuildID:   buildID,
		Size:      size,
		CreatedAt: time.Now(),
	}
}
