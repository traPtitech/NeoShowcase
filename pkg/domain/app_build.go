package domain

import (
	"slices"
	"time"

	"github.com/samber/lo"

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
	Artifacts     []*Artifact               // for static app
	RuntimeImage  optional.Of[RuntimeImage] // for runtime app if exists
}

func NewBuild(app *Application, env []*Environment) *Build {
	return &Build{
		ID:            NewID(),
		Commit:        app.Commit,
		ConfigHash:    app.Config.Hash(env),
		Status:        BuildStatusQueued,
		ApplicationID: app.ID,
		QueuedAt:      time.Now(),
	}
}

func (b *Build) GetWebsiteArtifact() (artifact *Artifact, ok bool) {
	slices.SortFunc(b.Artifacts, ds.MoreFunc(func(a *Artifact) int64 { return a.CreatedAt.UnixNano() }))
	return lo.Find(b.Artifacts, func(a *Artifact) bool { return a.Name == BuilderStaticArtifactName })
}
