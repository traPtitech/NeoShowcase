package domain

import (
	"context"
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
	Status        BuildStatus
	ApplicationID string
	QueuedAt      time.Time
	StartedAt     optional.Of[time.Time]
	UpdatedAt     optional.Of[time.Time]
	FinishedAt    optional.Of[time.Time]
	Retriable     bool
	Artifacts     []*Artifact
}

func NewBuild(applicationID string, commit string) *Build {
	return &Build{
		ID:            NewID(),
		Commit:        commit,
		Status:        BuildStatusQueued,
		ApplicationID: applicationID,
		QueuedAt:      time.Now(),
	}
}

func (b *Build) GetWebsiteArtifact() (artifact *Artifact, ok bool) {
	slices.SortFunc(b.Artifacts, ds.MoreFunc(func(a *Artifact) int64 { return a.CreatedAt.UnixNano() }))
	return lo.Find(b.Artifacts, func(a *Artifact) bool { return a.Name == BuilderStaticArtifactName })
}

// GetSuccessBuilds returns a map of (app id + build id) -> build.
func GetSuccessBuilds(ctx context.Context, buildRepo BuildRepository, apps []*Application) (map[string]*Build, error) {
	commits := make(map[string]struct{}, 2*len(apps))
	for _, app := range apps {
		commits[app.WantCommit] = struct{}{}
		commits[app.CurrentCommit] = struct{}{}
	}
	builds, err := buildRepo.GetBuilds(ctx, GetBuildCondition{CommitIn: optional.From(lo.Keys(commits)), Status: optional.From(BuildStatusSucceeded)})
	if err != nil {
		return nil, err
	}
	// Last succeeded builds for each app+commit
	slices.SortFunc(builds, func(a, b *Build) bool { return a.QueuedAt.Before(b.QueuedAt) })
	buildMap := lo.SliceToMap(builds, func(b *Build) (string, *Build) { return b.ApplicationID + b.Commit, b })
	return buildMap, nil
}
