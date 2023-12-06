package giteaintegration

import (
	"context"
	"fmt"
	"time"

	"code.gitea.io/sdk/gitea"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
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
	cancel      func()

	gitRepo  domain.GitRepositoryRepository
	appRepo  domain.ApplicationRepository
	userRepo domain.UserRepository
}

func NewIntegration(
	c Config,
	gitRepo domain.GitRepositoryRepository,
	appRepo domain.ApplicationRepository,
	userRepo domain.UserRepository,
) (*Integration, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	client, err := gitea.NewClient(
		c.URL,
		gitea.SetToken(c.Token),
		// Skip version check (traP fork versioned such as 'traP-1.19.3-2' regarded "Malformed version")
		gitea.SetGiteaVersion(""),
	)
	if err != nil {
		return nil, err
	}
	return &Integration{
		c:           client,
		interval:    time.Duration(c.IntervalSeconds) * time.Second,
		concurrency: c.Concurrency,

		gitRepo:  gitRepo,
		appRepo:  appRepo,
		userRepo: userRepo,
	}, nil
}

func (i *Integration) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel
	go loop.Loop(ctx, i.sync, i.interval, true)
	return nil
}

func (i *Integration) Shutdown() error {
	if i.cancel != nil {
		i.cancel()
	}
	return nil
}
