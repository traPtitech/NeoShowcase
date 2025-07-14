package domain

import (
	"time"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

const BuilderStaticArtifactName = "website.tar.gz"

const BuilderFunctionArtifactName = "function.js"

type Artifact struct {
	ID        string
	Name      string
	BuildID   string
	Size      int64
	CreatedAt time.Time
	DeletedAt optional.Of[time.Time]
}

func NewArtifact(buildID string, name string, size int64) *Artifact {
	return &Artifact{
		ID:        NewID(),
		Name:      name,
		BuildID:   buildID,
		Size:      size,
		CreatedAt: time.Now(),
	}
}
