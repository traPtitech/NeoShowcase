package dockerimpl

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/traPtitech/neoshowcase/pkg/container"
)

func (m *Manager) List(ctx context.Context) (*container.ListResult, error) {
	containers, err := m.c.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": {fmt.Sprintf("%s=true", appContainerLabel)}},
		Context: ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	var result []container.Container
	for _, apiContainers := range containers {
		result = append(result, container.Container{
			ApplicationID: apiContainers.Labels[appContainerApplicationIDLabel],
			EnvironmentID: apiContainers.Labels[appContainerEnvironmentIDLabel],
			State:         getContainerState(apiContainers.State),
		})
	}

	return &container.ListResult{
		Containers: result,
	}, nil
}

func getContainerState(state string) container.State {
	switch state {
	case "Created":
		return container.StateStopped
	case "Restarting":
		return container.StateRestarting
	case "Running":
		return container.StateRunning
	case "Paused":
		return container.StateOther
	case "Exited":
		return container.StateStopped
	case "Dead":
		return container.StateOther
	default:
		return container.StateOther
	}
}
