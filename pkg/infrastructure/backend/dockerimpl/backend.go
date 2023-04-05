package dockerimpl

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
)

type IngressConfDirPath string

const (
	appNetwork          = "neoshowcase_apps"
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
	c              *docker.Client
	bus            domain.Bus
	ingressConfDir string
	ssURL          string

	dockerEvent chan *docker.APIEvents
	reloadLock  sync.Mutex
}

func NewDockerBackend(
	c *docker.Client,
	bus domain.Bus,
	path IngressConfDirPath,
	ss domain.StaticServerConnectivityConfig,
) domain.Backend {
	return &dockerBackend{
		c:              c,
		bus:            bus,
		ingressConfDir: string(path),
		ssURL:          ss.URL,
	}
}

func (b *dockerBackend) Start(_ context.Context) error {
	// showcase用のネットワークを用意
	if err := initNetworks(b.c); err != nil {
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
			switch ev.Action {
			case "start":
				if ev.Actor.Attributes[appLabel] == "true" {
					b.bus.Publish(event.ContainerAppStarted, domain.Fields{
						"application_id": ev.Actor.Attributes[appIDLabel],
					})
				}
			case "stop":
				if ev.Actor.Attributes[appLabel] == "true" {
					b.bus.Publish(event.ContainerAppStopped, domain.Fields{
						"application_id": ev.Actor.Attributes[appIDLabel],
					})
				}
			}
		}
	}
}

func (b *dockerBackend) Dispose(_ context.Context) error {
	_ = b.c.RemoveEventListener(b.dockerEvent)
	close(b.dockerEvent)
	return nil
}

func initNetworks(c *docker.Client) error {
	networks, err := c.ListNetworks()
	if err != nil {
		return errors.Wrap(err, "failed to list networks")
	}
	for _, network := range networks {
		if network.Name == appNetwork {
			return nil
		}
	}

	_, err = c.CreateNetwork(docker.CreateNetworkOptions{
		Name: appNetwork,
	})
	return err
}

func containerLabels(app *domain.Application) map[string]string {
	return map[string]string{
		appLabel:            "true",
		appIDLabel:          app.ID,
		appRestartedAtLabel: app.UpdatedAt.Format(time.RFC3339),
	}
}

func containerName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

func networkName(appID string) string {
	return fmt.Sprintf("%s.nsapp.internal", appID)
}

func traefikName(website *domain.Website) string {
	s := fmt.Sprintf("nsapp-%s%s",
		strings.ReplaceAll(website.FQDN, ".", "-"),
		strings.ReplaceAll(website.PathPrefix, "/", "-"),
	)
	return strings.TrimSuffix(s, "-")
}

func stripMiddlewareName(website *domain.Website) string {
	return traefikName(website) + "-strip"
}

func ssHeaderMiddlewareName(ss *domain.StaticSite) string {
	return fmt.Sprintf("nsapp-ss-header-%s", ss.Application.ID)
}
