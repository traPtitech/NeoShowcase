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
			State:         apiContainers.State, // TODO
		})
	}

	return &container.ListResult{
		Containers: result,
	}, nil
}
