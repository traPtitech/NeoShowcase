package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/container"
)

const (
	appNetwork                     = "neoshowcase_apps"
	appContainerLabel              = "neoshowcase.trap.jp/app"
	appContainerApplicationIDLabel = "neoshowcase.trap.jp/appId"
	timeout                        = 5
)

type Manager struct {
	c           *docker.Client
	bus         *hub.Hub
	dockerEvent chan *docker.APIEvents
}

func NewManager(eventbus *hub.Hub) (*Manager, error) {
	// Dockerデーモンに接続 (DooD)
	d, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to connect docker daemon: %w", err)
	}

	// showcase用のネットワークを用意
	if err := initNetworks(d); err != nil {
		return nil, fmt.Errorf("failed to init networks: %w", err)
	}

	m := &Manager{
		c:           d,
		bus:         eventbus,
		dockerEvent: make(chan *docker.APIEvents, 10),
	}

	if err := d.AddEventListener(m.dockerEvent); err != nil {
		close(m.dockerEvent)
		return nil, fmt.Errorf("failed to add event listener: %w", err)
	}
	go m.eventListener()

	return m, nil
}

func (m *Manager) eventListener() {
	for ev := range m.dockerEvent {
		log.Debug(ev)
		// https://docs.docker.com/engine/reference/commandline/events/
		switch ev.Type {
		case "container":
			switch ev.Action {
			case "start":
				appID, ok := ev.Actor.Attributes[appContainerApplicationIDLabel]
				if ok {
					m.bus.Publish(hub.Message{
						Name: container.EventAppStarted,
						Fields: map[string]interface{}{
							"appId": appID,
						},
					})
				}
			case "stop":
				appID, ok := ev.Actor.Attributes[appContainerApplicationIDLabel]
				if ok {
					m.bus.Publish(hub.Message{
						Name: container.EventAppStopped,
						Fields: map[string]interface{}{
							"appId": appID,
						},
					})
				}
			}
		}
	}
}

func (m *Manager) Dispose(ctx context.Context) error {
	_ = m.c.RemoveEventListener(m.dockerEvent)
	close(m.dockerEvent)
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

func containerName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}
