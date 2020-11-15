package dockerimpl

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/leandro-lugaresi/hub"
)

const (
	appNetwork                     = "neoshowcase_apps"
	appContainerLabel              = "neoshowcase.trap.jp/app"
	appContainerApplicationIDLabel = "neoshowcase.trap.jp/appId"
)

type Manager struct {
	c   *docker.Client
	bus *hub.Hub
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

	return &Manager{
		c:   d,
		bus: eventbus,
	}, nil
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
