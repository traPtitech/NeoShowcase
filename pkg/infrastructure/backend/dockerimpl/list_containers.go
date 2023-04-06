package dockerimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/samber/lo"

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
		return nil, errors.Wrap(err, "failed to fetch containers")
	}

	if len(containers) == 0 {
		return &domain.Container{
			ApplicationID: appID,
			State:         domain.ContainerStateMissing,
		}, nil
	}
	return &domain.Container{
		ApplicationID: appID,
		State:         getContainerState(&containers[0]),
	}, nil
}

func (b *dockerBackend) ListContainers(ctx context.Context) ([]*domain.Container, error) {
	containers, err := b.c.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": {fmt.Sprintf("%s=true", appLabel)}},
		Context: ctx,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch containers")
	}

	result := lo.Map(containers, func(c docker.APIContainers, i int) *domain.Container {
		return &domain.Container{
			ApplicationID: c.Labels[appIDLabel],
			State:         getContainerState(&c),
		}
	})
	return result, nil
}

func getContainerState(c *docker.APIContainers) domain.ContainerState {
	// https://docs.docker.com/engine/api/v1.42/#tag/Container/operation/ContainerList
	switch strings.ToLower(c.State) {
	case "created":
		return domain.ContainerStateStarting
	case "restarting":
		return domain.ContainerStateRunning // to match with k8s pod phase
	case "running":
		return domain.ContainerStateRunning
	case "exited":
		status := strings.ToLower(c.Status)
		if strings.HasPrefix(status, "exited (0)") || strings.HasPrefix(status, "exit 0") {
			return domain.ContainerStateExited
		} else {
			return domain.ContainerStateErrored
		}
	case "dead":
		return domain.ContainerStateErrored
	default:
		return domain.ContainerStateUnknown
	}
}
