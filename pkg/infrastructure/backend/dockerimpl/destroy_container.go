package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	for _, website := range app.Websites {
		err := b.unregisterIngress(ctx, app, website)
		if err != nil {
			return fmt.Errorf("failed to unregister ingress: %w", err)
		}
	}
	err := b.c.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerName(app.ID),
		RemoveVolumes: true,
		Force:         true,
		Context:       ctx,
	})
	if err != nil {
		return fmt.Errorf("failed to destroy container: %w", err)
	}
	return nil
}
