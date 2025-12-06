package dockerimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (b *Backend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
	containers, err := b.c.ContainerList(ctx, client.ContainerListOptions{
		All: true,
		Filters: make(client.Filters).
			Add("label", fmt.Sprintf("%s=true", appLabel), fmt.Sprintf("%s=%s", appIDLabel, appID)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch containers")
	}

	if len(containers.Items) == 0 {
		return &domain.Container{
			ApplicationID: appID,
			State:         domain.ContainerStateMissing,
		}, nil
	}
	state, msg := getContainerState(&containers.Items[0])
	return &domain.Container{
		ApplicationID: appID,
		State:         state,
		Message:       msg,
	}, nil
}

func (b *Backend) ListContainers(ctx context.Context) ([]*domain.Container, error) {
	containers, err := b.c.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: make(client.Filters).Add("label", fmt.Sprintf("%s=true", appLabel)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch containers")
	}

	result := ds.Map(containers.Items, func(c container.Summary) *domain.Container {
		state, msg := getContainerState(&c)
		return &domain.Container{
			ApplicationID: c.Labels[appIDLabel],
			State:         state,
			Message:       msg,
		}
	})
	return result, nil
}

func getContainerState(c *container.Summary) (state domain.ContainerState, message string) {
	// https://docs.docker.com/engine/api/v1.42/#tag/Container/operation/ContainerList
	switch c.State {
	case container.StateCreated:
		return domain.ContainerStateStarting, c.Status
	case container.StateRestarting:
		return domain.ContainerStateRestarting, c.Status
	case container.StateRunning:
		return domain.ContainerStateRunning, c.Status
	case container.StateExited:
		status := strings.ToLower(c.Status)
		if strings.HasPrefix(status, "exited (0)") || strings.HasPrefix(status, "exit 0") {
			return domain.ContainerStateExited, c.Status
		} else {
			return domain.ContainerStateErrored, c.Status
		}
	case container.StateDead:
		return domain.ContainerStateErrored, c.Status
	default:
		return domain.ContainerStateUnknown, c.Status
	}
}
