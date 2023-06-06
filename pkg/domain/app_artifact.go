package domain

import (
	"time"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type Artifact struct {
	ID        string
	BuildID   string
	Size      int64
	CreatedAt time.Time
	DeletedAt optional.Of[time.Time]
}

func NewArtifact(buildID string, size int64) *Artifact {
	return &Artifact{
		ID:        NewID(),
		BuildID:   buildID,
		Size:      size,
		CreatedAt: time.Now(),
	}
}
