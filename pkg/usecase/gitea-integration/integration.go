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
	URL             string `mapstructure:"url" yaml:"url"`
	Token           string `mapstructure:"token" yaml:"token"`
	IntervalSeconds int    `mapstructure:"intervalSeconds" yaml:"intervalSeconds"`
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("provide admin gitea token")
	}
	if c.IntervalSeconds <= 0 {
		return errors.New("provide positive interval seconds")
	}
	return nil
}

type Integration struct {
	c        *gitea.Client
	interval time.Duration
	cancel   func()

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

	client, err := gitea.NewClient(c.URL, gitea.SetToken(c.Token))
	if err != nil {
		return nil, err
	}
	return &Integration{
		c:        client,
		interval: time.Duration(c.IntervalSeconds) * time.Second,

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
