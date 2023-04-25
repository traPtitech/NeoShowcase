package dockerimpl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	appLabel            = "neoshowcase.trap.jp/app"
	appIDLabel          = "neoshowcase.trap.jp/appId"
	appRestartedAtLabel = "neoshowcase.trap.jp/restartedAt"
)

const (
	traefikRuntimeFilename = "apps.yaml"
	traefikSSFilename      = "ss.yaml"
	traefikSSServiceName   = "ss"
)

type dockerBackend struct {
	c         *docker.Client
	conf      Config
	image     builder.ImageConfig
	eventSubs domain.PubSub[*domain.ContainerEvent]

	dockerEvent chan *docker.APIEvents
	subsLock    sync.Mutex

	reloadLock sync.Mutex
}

func NewDockerBackend(
	c *docker.Client,
	conf Config,
	image builder.ImageConfig,
) (domain.Backend, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	return &dockerBackend{
		c:     c,
		conf:  conf,
		image: image,
	}, nil
}

func (b *dockerBackend) Start(_ context.Context) error {
	// showcase用のネットワークを用意
	if err := b.initNetworks(); err != nil {
		return errors.Wrap(err, "failed to init networks")
	}

	b.dockerEvent = make(chan *docker.APIEvents, 10)
	if err := b.c.AddEventListener(b.dockerEvent); err != nil {
		close(b.dockerEvent)
		return errors.Wrap(err, "failed to add event listener")
	}
	go b.eventListener()

	return nil
}

func (b *dockerBackend) eventListener() {
	for ev := range b.dockerEvent {
		// https://docs.docker.com/engine/reference/commandline/events/
		switch ev.Type {
		case "container":
			appID, ok := ev.Actor.Attributes[appIDLabel]
			if !ok {
				continue
			}
			b.eventSubs.Publish(&domain.ContainerEvent{ApplicationID: appID})
		}
	}
}

func (b *dockerBackend) Dispose(_ context.Context) error {
	_ = b.c.RemoveEventListener(b.dockerEvent)
	close(b.dockerEvent)
	return nil
}

func (b *dockerBackend) AuthAllowed(fqdn string) bool {
	for _, ac := range b.conf.Middlewares.Auth {
		if domain.MatchDomain(ac.Domain, fqdn) {
			return true
		}
	}
	return false
}

func (b *dockerBackend) targetAuth(fqdn string) *authConf {
	for _, ac := range b.conf.Middlewares.Auth {
		if domain.MatchDomain(ac.Domain, fqdn) {
			return ac
		}
	}
	return nil
}

func (b *dockerBackend) ListenContainerEvents() (sub <-chan *domain.ContainerEvent, unsub func()) {
	return b.eventSubs.Subscribe()
}

func (b *dockerBackend) initNetworks() error {
	networks, err := b.c.ListNetworks()
	if err != nil {
		return errors.Wrap(err, "failed to list networks")
	}
	for _, network := range networks {
		if network.Name == b.conf.Network {
			return nil
		}
	}

	_, err = b.c.CreateNetwork(docker.CreateNetworkOptions{
		Name: b.conf.Network,
	})
	return err
}

func (b *dockerBackend) authConfig() docker.AuthConfiguration {
	if b.image.Registry.Username == "" && b.image.Registry.Password == "" {
		return docker.AuthConfiguration{}
	}
	return docker.AuthConfiguration{
		Username: b.image.Registry.Username,
		Password: b.image.Registry.Password,
	}
}

func (b *dockerBackend) containerLabels(app *domain.Application) map[string]string {
	return ds.MergeMap(b.conf.labels(), map[string]string{
		appLabel:            "true",
		appIDLabel:          app.ID,
		appRestartedAtLabel: app.UpdatedAt.Format(time.RFC3339),
	})
}

func containerName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

func networkName(appID string) string {
	return fmt.Sprintf("%s.nsapp.internal", appID)
}

func traefikName(website *domain.Website) string {
	return fmt.Sprintf("nsapp-%s", website.ID)
}

func stripMiddlewareName(website *domain.Website) string {
	return traefikName(website) + "-strip"
}

func ssHeaderMiddlewareName(ss *domain.StaticSite) string {
	return fmt.Sprintf("nsapp-ss-header-%s", ss.Application.ID)
}
