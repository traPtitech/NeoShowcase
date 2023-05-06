package webhook

import (
	"context"
	"fmt"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ReceiverConfig struct {
	BasePath string `mapstructure:"basePath" yaml:"basePath"`
	Port     int    `mapstructure:"port" yaml:"port"`
}

type Receiver struct {
	config  ReceiverConfig
	gitRepo domain.GitRepositoryRepository
	fetcher usecase.RepositoryFetcherService

	echo *echo.Echo
}

func NewReceiver(
	config ReceiverConfig,
	gitRepo domain.GitRepositoryRepository,
	fetcher usecase.RepositoryFetcherService,
) *Receiver {
	r := &Receiver{
		config:  config,
		gitRepo: gitRepo,
		fetcher: fetcher,
	}

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Any(path.Join(config.BasePath, "github"), r.githubHandler)
	e.Any(path.Join(config.BasePath, "gitea"), r.giteaHandler)
	r.echo = e

	return r
}

func (r *Receiver) Start(_ context.Context) error {
	return r.echo.Start(fmt.Sprintf(":%d", r.config.Port))
}

func (r *Receiver) Shutdown(ctx context.Context) error {
	return r.echo.Shutdown(ctx)
}

func (r *Receiver) updateURLs(urls []string) {
	ctx := context.Background()
	repos, err := r.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{URLs: optional.From(urls)})
	if err != nil {
		log.Errorf("getting repositories by url: %+v", err)
		return
	}
	if len(repos) == 0 {
		return
	}
	ids := lo.Map(repos, func(repo *domain.Repository, _ int) string { return repo.ID })
	r.fetcher.Fetch(ids)
}
