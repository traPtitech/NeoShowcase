package webhook

import (
	"context"
	"fmt"
	"path"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ReceiverConfig struct {
	BasePath string `mapstructure:"basePath" yaml:"basePath"`
	Port     int    `mapstructure:"port" yaml:"port"`
}

// Receiver is responsible for receiving webhook events and notifying repository refresh.
//
// We should NOT trust the payload of webhooks, as they could be crafted and there is no secret.
// This behavior is inspired from ArgoCD Webhook behavior.
//
// Cited from https://argo-cd.readthedocs.io/en/stable/operator-manual/webhook/
//
//	Configuring a webhook shared secret is optional, since Argo CD will still refresh applications
//	related to the Git repository, even with unauthenticated webhook events.
//	This is safe to do since the contents of webhook payloads are considered untrusted,
//	and will only result in a refresh of the application (a process which already occurs at three-minute intervals).
type Receiver struct {
	config           ReceiverConfig
	gitRepo          domain.GitRepositoryRepository
	fetcher          repofetcher.Service
	giteaIntegration domain.ControllerGiteaIntegrationService

	echo *echo.Echo
}

func NewReceiver(
	config ReceiverConfig,
	gitRepo domain.GitRepositoryRepository,
	fetcher repofetcher.Service,
	giteaIntegration domain.ControllerGiteaIntegrationService,
) *Receiver {
	r := &Receiver{
		config:           config,
		gitRepo:          gitRepo,
		fetcher:          fetcher,
		giteaIntegration: giteaIntegration,
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

// updateURLs notifies repositories to refresh.
func (r *Receiver) updateURLs(urls []string) {
	ctx := context.Background()
	repos, err := r.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{URLs: optional.From(urls)})
	if err != nil {
		log.Errorf("getting repositories by url: %+v", err)
		return
	}
	for _, repo := range repos {
		r.fetcher.Fetch(repo.ID)
	}
}
