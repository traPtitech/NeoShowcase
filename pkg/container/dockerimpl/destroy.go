package dockerimpl

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/traPtitech/neoshowcase/pkg/container"
)

func (m *Manager) Destroy(ctx context.Context, args container.DestroyArgs) (*container.DestroyResult, error) {
	err := m.c.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerName(args.ApplicationID),
		RemoveVolumes: true,
		Force:         true,
		Context:       ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to destroy container: %w", err)
	}
	return &container.DestroyResult{}, nil
}
