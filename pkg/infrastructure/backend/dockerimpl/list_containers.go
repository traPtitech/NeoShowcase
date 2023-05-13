package dockerimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (b *dockerBackend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
	containers, err := b.c.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("%s=true", appLabel)),
			filters.Arg("label", fmt.Sprintf("%s=%s", appIDLabel, appID)),
		),
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
	containers, err := b.c.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("%s=true", appLabel)),
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch containers")
	}

	result := ds.Map(containers, func(c types.Container) *domain.Container {
		return &domain.Container{
			ApplicationID: c.Labels[appIDLabel],
			State:         getContainerState(&c),
		}
	})
	return result, nil
}

func getContainerState(c *types.Container) domain.ContainerState {
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
