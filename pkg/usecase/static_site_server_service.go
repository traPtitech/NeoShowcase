package usecase

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type SiteReloadTarget struct {
	ApplicationID string
	BuildID       string
}

type StaticSiteServerService interface {
	Reload(ctx context.Context) error
}

type staticSiteServerService struct {
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	engine    domain.Engine
}

func NewStaticSiteServerService(
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	engine domain.Engine,
) StaticSiteServerService {
	return &staticSiteServerService{
		appRepo:   appRepo,
		buildRepo: buildRepo,
		engine:    engine,
	}
}

func (s *staticSiteServerService) Reload(ctx context.Context) error {
	data, err2 := domain.GetActiveWebsites(ctx, s.appRepo, s.buildRepo)
	if err2 != nil {
		return err2
	}

	return s.engine.Reconcile(data)
}
