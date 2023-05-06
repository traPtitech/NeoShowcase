package usecase

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

type SiteReloadTarget struct {
	ApplicationID string
	BuildID       string
}

type StaticSiteServerService interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type staticSiteServerService struct {
	client    domain.ControllerSSGenServiceClient
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	engine    domain.SSEngine

	cancel     func()
	reloadLock sync.Mutex
}

func NewStaticSiteServerService(
	client domain.ControllerSSGenServiceClient,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	engine domain.SSEngine,
) StaticSiteServerService {
	return &staticSiteServerService{
		client:    client,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		engine:    engine,
	}
}

func (s *staticSiteServerService) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	go retry.Do(ctx, func(ctx context.Context) error {
		return s.client.ConnectSSGen(ctx, s.onRequest)
	}, 1*time.Second, 60*time.Second)

	go func() {
		err := s.reload(context.Background())
		if err != nil {
			log.Errorf("failed to perform initial reload: %+v", err)
		}
	}()

	return nil
}

func (s *staticSiteServerService) Shutdown(_ context.Context) error {
	s.cancel()
	return nil
}

func (s *staticSiteServerService) onRequest(req *pb.SSGenRequest) {
	switch req.Type {
	case pb.SSGenRequest_RELOAD:
		err := s.reload(context.Background())
		if err != nil {
			log.Errorf("failed to reload static server: %+v", err)
		}
	}
}

func (s *staticSiteServerService) reload(ctx context.Context) error {
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
