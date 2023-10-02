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

func (b *Backend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
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
	state, msg := getContainerState(&containers[0])
	return &domain.Container{
		ApplicationID: appID,
		State:         state,
		Message:       msg,
	}, nil
}

func (b *Backend) ListContainers(ctx context.Context) ([]*domain.Container, error) {
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
		state, msg := getContainerState(&c)
		return &domain.Container{
			ApplicationID: c.Labels[appIDLabel],
			State:         state,
			Message:       msg,
		}
	})
	return result, nil
}

func getContainerState(c *types.Container) (state domain.ContainerState, message string) {
	// https://docs.docker.com/engine/api/v1.42/#tag/Container/operation/ContainerList
	switch strings.ToLower(c.State) {
	case "created":
		return domain.ContainerStateStarting, c.Status
	case "restarting":
		return domain.ContainerStateRestarting, c.Status
	case "running":
		return domain.ContainerStateRunning, c.Status
	case "exited":
		status := strings.ToLower(c.Status)
		if strings.HasPrefix(status, "exited (0)") || strings.HasPrefix(status, "exit 0") {
			return domain.ContainerStateExited, c.Status
		} else {
			return domain.ContainerStateErrored, c.Status
		}
	case "dead":
		return domain.ContainerStateErrored, c.Status
	default:
		return domain.ContainerStateUnknown, c.Status
	}
}
