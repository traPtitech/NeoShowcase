package giteaintegration

import (
	"context"
	"errors"
	"time"

	"code.gitea.io/sdk/gitea"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
)

type Config struct {
	URL             string
	Token           string
	IntervalSeconds int
	ListIntervalMs  int
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("provide admin gitea token")
	}
	if c.IntervalSeconds <= 0 {
		return errors.New("provide positive interval seconds")
	}
	if c.ListIntervalMs <= 0 {
		return errors.New("provide positive list interval ms")
	}
	return nil
}

type Integration struct {
	c            *gitea.Client
	interval     time.Duration
	listInterval time.Duration
	cancel       func()

	gitRepo  domain.GitRepositoryRepository
	userRepo domain.UserRepository
}

func NewIntegration(
	c Config,
	gitRepo domain.GitRepositoryRepository,
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
		c:            client,
		interval:     time.Duration(c.IntervalSeconds) * time.Second,
		listInterval: time.Duration(c.ListIntervalMs) * time.Millisecond,

		gitRepo:  gitRepo,
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
