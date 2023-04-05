package dockerimpl

import (
	"context"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func TestDockerBackend_CreateContainer(t *testing.T) {
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("存在しないイメージを指定", func(t *testing.T) {
		t.Parallel()
		app := domain.Application{
			ID: "test",
		}
		err := m.Synchronize(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: "not-found",
			ImageTag:  "latest",
		}})
		assert.Error(t, err)
	})

	t.Run("コンテナを正常に作成して起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"

		app := domain.Application{
			ID: appID,
		}
		err := m.Synchronize(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
		}})
		require.NoError(t, err)

		cont, err := c.InspectContainerWithOptions(docker.InspectContainerOptions{
			ID: containerName(appID),
		})
		require.NoError(t, err)

		assert.Equal(t, cont.Config.Image, image+":latest")
		assert.Equal(t, cont.Config.Labels[appLabel], "true")
		assert.Equal(t, cont.Config.Labels[appIDLabel], appID)

		require.NoError(t, c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            cont.ID,
			RemoveVolumes: true,
			Force:         true,
		}))
	})

	t.Run("コンテナを正常に作成して起動 (Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pij0bij90j20"

		app := domain.Application{
			ID: appID,
		}
		err := m.Synchronize(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
		}})
		require.NoError(t, err)

		err = m.Synchronize(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
			Restart:   true,
		}})
		require.NoError(t, err)

		cont, err := c.InspectContainerWithOptions(docker.InspectContainerOptions{
			ID: containerName(appID),
		})
		require.NoError(t, err)

		assert.Equal(t, cont.Config.Image, image+":latest")
		assert.Equal(t, cont.Config.Labels[appLabel], "true")
		assert.Equal(t, cont.Config.Labels[appIDLabel], appID)

		require.NoError(t, c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            cont.ID,
			RemoveVolumes: true,
			Force:         true,
		}))
	})
}
