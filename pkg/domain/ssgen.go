package domain

import (
	"context"

	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type StaticServerDocumentRootPath string

type StaticServer interface {
	Start(ctx context.Context) error
	Reconcile(sites []*StaticSite) error
	Shutdown(ctx context.Context) error
}

type StaticSite struct {
	Application *Application
	Website     *Website
	ArtifactID  string
	SPA         bool
}

func GetActiveStaticSites(
	ctx context.Context,
	cluster *discovery.Cluster,
	appRepo ApplicationRepository,
	buildRepo BuildRepository,
) ([]*StaticSite, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		DeployType: optional.From(DeployTypeStatic),
		Running:    optional.From(true),
	})
	if err != nil {
		return nil, err
	}
	// Shard by app ID
	applications = lo.Filter(applications, func(app *Application, _ int) bool {
		return cluster.IsAssigned(app.ID)
	})

	buildIDs := lo.FilterMap(applications, func(app *Application, _ int) (string, bool) {
		return app.CurrentBuild, app.CurrentBuild != ""
	})
	builds, err := buildRepo.GetBuilds(ctx, GetBuildCondition{IDIn: optional.From(buildIDs)})
	if err != nil {
		return nil, err
	}
	buildsMap := lo.SliceToMap(builds, func(b *Build) (string, *Build) { return b.ID, b })
	var sites []*StaticSite
	for _, app := range applications {
		build, ok := buildsMap[app.CurrentBuild]
		if !ok {
			continue
		}
		artifact, ok := build.GetWebsiteArtifact()
		if !ok {
			continue
		}
		for _, website := range app.Websites {
			sites = append(sites, &StaticSite{
				Application: app,
				Website:     website,
				ArtifactID:  artifact.ID,
				SPA:         app.Config.BuildConfig.GetStaticConfig().SPA,
			})
		}
	}
	return sites, nil
}
