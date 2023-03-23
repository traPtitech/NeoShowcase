package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
	containers, err := b.c.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{"label": {
			fmt.Sprintf("%s=true", appLabel),
			fmt.Sprintf("%s=%s", appIDLabel, appID),
		}},
		Context: ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}
	if len(containers) == 0 {
		return nil, domain.ErrContainerNotFound
	}

	apiContainer := containers[0]
	return &domain.Container{
		ApplicationID: appID,
		State:         getContainerState(apiContainer.State),
	}, nil
}

func (b *dockerBackend) ListContainers(ctx context.Context) ([]domain.Container, error) {
	containers, err := b.c.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": {fmt.Sprintf("%s=true", appLabel)}},
		Context: ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	var result []domain.Container
	for _, apiContainers := range containers {
		result = append(result, domain.Container{
			ApplicationID: apiContainers.Labels[appIDLabel],
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
