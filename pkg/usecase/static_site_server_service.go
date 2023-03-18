package usecase

import (
	"context"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type SiteReloadTarget struct {
	ApplicationID string
	BuildID       string
}

type StaticSiteServerService interface {
	Reload(ctx context.Context) error
}

type staticSiteServerService struct {
	appRepo   repository.ApplicationRepository
	buildRepo repository.BuildRepository
	engine    domain.Engine
}

func NewStaticSiteServerService(
	appRepo repository.ApplicationRepository,
	buildRepo repository.BuildRepository,
	engine domain.Engine,
) StaticSiteServerService {
	return &staticSiteServerService{
		appRepo:   appRepo,
		buildRepo: buildRepo,
		engine:    engine,
	}
}

func (s *staticSiteServerService) Reload(ctx context.Context) error {
	applications, err := s.appRepo.GetApplications(ctx, repository.GetApplicationCondition{
		BuildType: optional.From(builder.BuildTypeStatic),
		State:     optional.From(domain.ApplicationStateRunning),
	})
	if err != nil {
		return err
	}

	commits := lo.Map(applications, func(app *domain.Application, i int) string { return app.CurrentCommit })
	builds, err := s.buildRepo.GetBuildsInCommit(ctx, commits)
	if err != nil {
		return err
	}

	// Last succeeded builds for each commit
	builds = lo.Filter(builds, func(build *domain.Build, i int) bool { return build.Status == builder.BuildStatusSucceeded })
	slices.SortFunc(builds, func(a, b *domain.Build) bool { return a.StartedAt.Before(b.StartedAt) })
	commitToBuild := lo.SliceToMap(builds, func(b *domain.Build) (string, *domain.Build) { return b.Commit, b })

	var data []*domain.Site
	for _, app := range applications {
		build, ok := commitToBuild[app.CurrentCommit]
		if !ok {
			continue
		}
		if !build.Artifact.Valid {
			continue
		}
		for _, website := range app.Websites {
			data = append(data, &domain.Site{
				ID:            website.ID,
				FQDN:          website.FQDN,
				ArtifactID:    build.Artifact.V.ID,
				ApplicationID: app.ID,
			})
		}
	}

	return s.engine.Reconcile(data)
}
