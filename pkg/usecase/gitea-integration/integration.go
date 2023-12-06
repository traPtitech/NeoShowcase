package giteaintegration

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"code.gitea.io/sdk/gitea"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
	"github.com/traPtitech/neoshowcase/pkg/util/scutil"
)

type Config struct {
	URL             string
	Token           string
	IntervalSeconds int
	Concurrency     int
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("provide admin gitea token (got empty string)")
	}
	if c.IntervalSeconds <= 0 {
		return fmt.Errorf("provide positive interval seconds (got %v)", c.IntervalSeconds)
	}
	if c.Concurrency <= 0 {
		return fmt.Errorf("provide positive concurrency (got %v)", c.Concurrency)
	}
	return nil
}

type Integration struct {
	c           *gitea.Client
	interval    time.Duration
	concurrency int

	controller domain.ControllerGiteaIntegrationServiceClient
	gitRepo    domain.GitRepositoryRepository
	appRepo    domain.ApplicationRepository
	userRepo   domain.UserRepository

	cancel func()
	syncer *scutil.Coalescer
}

func NewIntegration(
	c Config,
	controller domain.ControllerGiteaIntegrationServiceClient,
	gitRepo domain.GitRepositoryRepository,
	appRepo domain.ApplicationRepository,
	userRepo domain.UserRepository,
) (*Integration, error) {
	// Validate config
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Generate gitea client
	client, err := gitea.NewClient(
		c.URL,
		gitea.SetToken(c.Token),
		// Skip version check (traP fork versioned such as 'traP-1.19.3-2' regarded "Malformed version")
		gitea.SetGiteaVersion(""),
	)
	if err != nil {
		return nil, err
	}

	i := &Integration{
		c:           client,
		interval:    time.Duration(c.IntervalSeconds) * time.Second,
		concurrency: c.Concurrency,

		controller: controller,
		gitRepo:    gitRepo,
		appRepo:    appRepo,
		userRepo:   userRepo,
	}
	i.syncer = scutil.NewCoalescer(i.syncAndLog)
	return i, nil
}

func (i *Integration) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel

	go retry.Do(ctx, func(ctx context.Context) error {
		return i.controller.Connect(ctx, i.onRequest(ctx))
	}, "connect to controller")
	go loop.Loop(ctx, func(ctx context.Context) {
		_ = i.syncer.Do(ctx)
	}, i.interval, true)

	return nil
}

func (i *Integration) onRequest(ctx context.Context) func(req *pb.GiteaIntegrationRequest) {
	return func(req *pb.GiteaIntegrationRequest) {
		switch req.Type {
		case pb.GiteaIntegrationRequest_RESYNC:
			go func() {
				_ = i.syncer.Do(ctx)
			}()
		}
	}
}

func (i *Integration) syncAndLog(ctx context.Context) error {
	err := i.sync(ctx)
	if err != nil {
		log.Errorf("failed to sync: %+v", err)
	}
	return nil
}

func (i *Integration) Shutdown() error {
	if i.cancel != nil {
		i.cancel()
	}
	return nil
}
