package domain

import (
	"context"

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
}

func GetActiveStaticSites(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*StaticSite, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		DeployType: optional.From(DeployTypeStatic),
		Running:    optional.From(true),
	})
	if err != nil {
		return nil, err
	}

	builds, err := GetSuccessBuilds(ctx, buildRepo, applications)
	if err != nil {
		return nil, err
	}
	var sites []*StaticSite
	for _, app := range applications {
		build, ok := builds[app.ID+app.CurrentCommit]
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
