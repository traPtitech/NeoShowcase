package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) ListContainers(ctx context.Context) ([]domain.Container, error) {
	containers, err := b.c.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": {fmt.Sprintf("%s=true", appContainerLabel)}},
		Context: ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	var result []domain.Container
	for _, apiContainers := range containers {
		result = append(result, domain.Container{
			ApplicationID: apiContainers.Labels[appContainerApplicationIDLabel],
			EnvironmentID: apiContainers.Labels[appContainerEnvironmentIDLabel],
			State:         getContainerState(apiContainers.State),
		})
	}
	return result, nil
}

func getContainerState(state string) domain.ContainerState {
	switch state {
	case "Created":
		return domain.ContainerStateStopped
	case "Restarting":
		return domain.ContainerStateRestarting
	case "Running":
		return domain.ContainerStateRunning
	case "Paused":
		return domain.ContainerStateOther
	case "Exited":
		return domain.ContainerStateStopped
	case "Dead":
		return domain.ContainerStateOther
	default:
		return domain.ContainerStateOther
	}
}
