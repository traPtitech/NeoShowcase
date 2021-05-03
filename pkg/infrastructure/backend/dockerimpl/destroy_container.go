package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func (b *dockerBackend) DestroyContainer(ctx context.Context, appID string, envID string) error {
	err := b.c.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerName(appID, envID),
		RemoveVolumes: true,
		Force:         true,
		Context:       ctx,
	})
	if err != nil {
		return fmt.Errorf("failed to destroy container: %w", err)
	}
	return nil
}
