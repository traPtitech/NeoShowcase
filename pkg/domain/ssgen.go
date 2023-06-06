package domain

import (
	"context"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type (
	StaticServerDocumentRootPath string
	StaticServerPort             int
)

type SSEngine interface {
	Start(ctx context.Context) error
	Reconcile(sites []*StaticSite) error
	Shutdown(ctx context.Context) error
}

type StaticSite struct {
	Application *Application
	Website     *Website
	ArtifactID  string
}

func GetActiveStaticSites(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*StaticSite, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		DeployType: optional.From(DeployTypeStatic),
		Running:    optional.From(true),
	})
	if err != nil {
		return nil, err
	}

	commits := ds.Map(applications, func(app *Application) string { return app.CurrentCommit })
	builds, err := buildRepo.GetBuilds(ctx, GetBuildCondition{CommitIn: optional.From(commits), Status: optional.From(BuildStatusSucceeded)})
	if err != nil {
		return nil, err
	}

	// Last succeeded builds for each app+commit
	slices.SortFunc(builds, func(a, b *Build) bool { return a.StartedAt.ValueOrZero().Before(b.StartedAt.ValueOrZero()) })
	buildMap := lo.SliceToMap(builds, func(b *Build) (string, *Build) { return b.ApplicationID + b.Commit, b })

	var sites []*StaticSite
	for _, app := range applications {
		build, ok := buildMap[app.ID+app.CurrentCommit]
		if !ok {
			continue
		}
		if !build.Artifact.Valid {
			continue
		}
		for _, website := range app.Websites {
			sites = append(sites, &StaticSite{
				Application: app,
				Website:     website,
				ArtifactID:  build.Artifact.V.ID,
			})
		}
	}
	return sites, nil
}
