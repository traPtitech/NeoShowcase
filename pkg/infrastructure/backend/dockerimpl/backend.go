package dockerimpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
)

type IngressConfDirPath string

const (
	appNetwork = "neoshowcase_apps"
	appLabel   = "neoshowcase.trap.jp/app"
	appIDLabel = "neoshowcase.trap.jp/appId"
)

const (
	traefikHTTPEntrypoint     = "web"
	traefikHTTPSEntrypoint    = "websecure"
	traefikAuthSoftMiddleware = "ns_auth_soft@file"
	traefikAuthHardMiddleware = "ns_auth_hard@file"
	traefikAuthMiddleware     = "ns_auth@file"
	traefikCertResolver       = "nsresolver@file"
	traefikSSFilename         = "ss.yaml"
	traefikSSServiceName      = "ss"
)

type dockerBackend struct {
	c              *docker.Client
	bus            domain.Bus
	ingressConfDir string

	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	ssURL     string

	dockerEvent chan *docker.APIEvents
	reloadLock  sync.Mutex
}

func NewDockerBackend(
	c *docker.Client,
	bus domain.Bus,
	path IngressConfDirPath,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	ss domain.StaticServerConnectivityConfig,
) domain.Backend {
	return &dockerBackend{
		c:              c,
		bus:            bus,
		ingressConfDir: string(path),

		appRepo:   appRepo,
		buildRepo: buildRepo,
		ssURL:     ss.URL,
	}
}

func (b *dockerBackend) Start(_ context.Context) error {
	// showcase用のネットワークを用意
	if err := initNetworks(b.c); err != nil {
		return fmt.Errorf("failed to init networks: %w", err)
	}

	b.dockerEvent = make(chan *docker.APIEvents, 10)
	if err := b.c.AddEventListener(b.dockerEvent); err != nil {
		close(b.dockerEvent)
		return fmt.Errorf("failed to add event listener: %w", err)
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
		return fmt.Errorf("failed to list networks: %w", err)
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

func containerLabels(appID string) map[string]string {
	return map[string]string{
		appLabel:   "true",
		appIDLabel: appID,
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

func configFilePrefix(app *domain.Application) string {
	return fmt.Sprintf("nsapp-%s-", app.ID)
}

func configFile(app *domain.Application, website *domain.Website) string {
	filename := configFilePrefix(app) +
		strings.ReplaceAll(website.FQDN, ".", "-") +
		strings.ReplaceAll(website.PathPrefix, "/", "-")
	filename = strings.TrimSuffix(filename, "-")
	filename += ".yaml"
	return filename
}
