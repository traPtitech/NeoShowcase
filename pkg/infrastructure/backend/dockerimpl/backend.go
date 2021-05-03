package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

const (
	appNetwork                     = "neoshowcase_apps"
	appContainerLabel              = "neoshowcase.trap.jp/app"
	appContainerApplicationIDLabel = "neoshowcase.trap.jp/appId"
	appContainerEnvironmentIDLabel = "neoshowcase.trap.jp/envId"
	timeout                        = 5
)

type dockerBackend struct {
	c           *docker.Client
	bus         eventbus.Bus
	dockerEvent chan *docker.APIEvents
}

func NewDockerBackend(c *docker.Client, bus eventbus.Bus) (backend.Backend, error) {
	// showcase用のネットワークを用意
	if err := initNetworks(c); err != nil {
		return nil, fmt.Errorf("failed to init networks: %w", err)
	}

	b := &dockerBackend{
		c:           c,
		bus:         bus,
		dockerEvent: make(chan *docker.APIEvents, 10),
	}

	if err := c.AddEventListener(b.dockerEvent); err != nil {
		close(b.dockerEvent)
		return nil, fmt.Errorf("failed to add event listener: %w", err)
	}
	go b.eventListener()

	return b, nil
}

func (b *dockerBackend) eventListener() {
	for ev := range b.dockerEvent {
		log.Debug(ev)
		// https://docs.docker.com/engine/reference/commandline/events/
		switch ev.Type {
		case "container":
			switch ev.Action {
			case "start":
				if ev.Actor.Attributes[appContainerLabel] == "true" {
					b.bus.Publish(event.ContainerAppStarted, eventbus.Fields{
						"application_id": ev.Actor.Attributes[appContainerApplicationIDLabel],
						"environment_id": ev.Actor.Attributes[appContainerEnvironmentIDLabel],
					})
				}
			case "stop":
				if ev.Actor.Attributes[appContainerLabel] == "true" {
					b.bus.Publish(event.ContainerAppStopped, eventbus.Fields{
						"application_id": ev.Actor.Attributes[appContainerApplicationIDLabel],
						"environment_id": ev.Actor.Attributes[appContainerEnvironmentIDLabel],
					})
				}
			}
		}
	}
}

func (b *dockerBackend) Dispose(ctx context.Context) error {
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

func containerName(appID, envID string) string {
	return fmt.Sprintf("nsapp-%s-%s", appID, envID)
}
