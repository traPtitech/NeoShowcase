package usecase

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

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
	engine    domain.SSEngine

	reloadLock sync.Mutex
}

func NewStaticSiteServerService(
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	engine domain.SSEngine,
) StaticSiteServerService {
	return &staticSiteServerService{
		appRepo:   appRepo,
		buildRepo: buildRepo,
		engine:    engine,
	}
}

func (s *staticSiteServerService) Reload(ctx context.Context) error {
	s.reloadLock.Lock()
	defer s.reloadLock.Unlock()

	start := time.Now()
	defer log.Infof("reloaded static server in %v", time.Since(start))

	data, err := domain.GetActiveStaticSites(ctx, s.appRepo, s.buildRepo)
	if err != nil {
		return err
	}
	return s.engine.Reconcile(data)
}
