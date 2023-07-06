package domain

import (
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type BuildStatus int

const (
	BuildStatusQueued BuildStatus = iota
	BuildStatusBuilding
	BuildStatusSucceeded
	BuildStatusFailed
	BuildStatusCanceled
	BuildStatusSkipped
)

func (t BuildStatus) IsFinished() bool {
	switch t {
	case BuildStatusSucceeded, BuildStatusFailed, BuildStatusCanceled, BuildStatusSkipped:
		return true
	default:
		return false
	}
}

type Build struct {
	ID            string
	Commit        string
	ConfigHash    string
	Status        BuildStatus
	ApplicationID string
	QueuedAt      time.Time
	StartedAt     optional.Of[time.Time]
	UpdatedAt     optional.Of[time.Time]
	FinishedAt    optional.Of[time.Time]
	Retriable     bool
	Artifacts     []*Artifact
}

func NewBuild(app *Application) *Build {
	return &Build{
		ID:            NewID(),
		Commit:        app.Commit,
		ConfigHash:    app.Config.Hash(),
		Status:        BuildStatusQueued,
		ApplicationID: app.ID,
		QueuedAt:      time.Now(),
	}
}

func (b *Build) GetWebsiteArtifact() (artifact *Artifact, ok bool) {
	slices.SortFunc(b.Artifacts, ds.MoreFunc(func(a *Artifact) int64 { return a.CreatedAt.UnixNano() }))
	return lo.Find(b.Artifacts, func(a *Artifact) bool { return a.Name == BuilderStaticArtifactName })
}
