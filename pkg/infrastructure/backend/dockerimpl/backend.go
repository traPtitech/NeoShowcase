package dockerimpl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type Config struct {
	ConfDir string `mapstructure:"confDir" yaml:"confDir"`
	SS      struct {
		URL string `mapstructure:"url" yaml:"url"`
	} `mapstructure:"ss" yaml:"ss"`
	Network string            `mapstructure:"network" yaml:"network"`
	Labels  map[string]string `mapstructure:"labels" yaml:"labels"`
	TLS     struct {
		CertResolver string `mapstructure:"certResolver" yaml:"certResolver"`
		Wildcard     struct {
			Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
		} `mapstructure:"wildcard" yaml:"wildcard"`
	} `mapstructure:"tls" yaml:"tls"`
}

func (c *Config) Validate() error {
	if err := c.TLS.Wildcard.Domains.Validate(); err != nil {
		return errors.Wrap(err, "docker.tls.wildcard.domains is invalid")
	}
	return nil
}

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
	c    *docker.Client
	bus  domain.Bus
	conf Config

	dockerEvent chan *docker.APIEvents
	reloadLock  sync.Mutex
}

func NewDockerBackend(
	c *docker.Client,
	bus domain.Bus,
	conf Config,
) (domain.Backend, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	return &dockerBackend{
		c:    c,
		bus:  bus,
		conf: conf,
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
			b.bus.Publish(event.AppContainerUpdated, domain.Fields{
				"application_id": appID,
			})
		}
	}
}

func (b *dockerBackend) Dispose(_ context.Context) error {
	_ = b.c.RemoveEventListener(b.dockerEvent)
	close(b.dockerEvent)
	return nil
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

func (b *dockerBackend) containerLabels(app *domain.Application) map[string]string {
	return ds.MergeMap(b.conf.Labels, map[string]string{
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
