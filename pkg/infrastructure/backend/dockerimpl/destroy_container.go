package dockerimpl

import (
	"context"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	err := b.destroyRuntimeIngresses(ctx, app)
	if err != nil {
		return nil
	}

	err = b.c.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerName(app.ID),
		RemoveVolumes: true,
		Force:         true,
		Context:       ctx,
	})
	if err != nil {
		return errors.Wrap(err, "failed to destroy container")
	}
	return nil
}
