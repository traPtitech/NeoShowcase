package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func (b *dockerBackend) DestroyContainer(ctx context.Context, appID string, branchID string) error {
	err := b.c.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerName(appID, branchID),
		RemoveVolumes: true,
		Force:         true,
		Context:       ctx,
	})
	if err != nil {
		return fmt.Errorf("failed to destroy container: %w", err)
	}
	return nil
}
