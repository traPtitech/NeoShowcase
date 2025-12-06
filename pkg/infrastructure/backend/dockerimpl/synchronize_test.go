package dockerimpl

import (
	"context"
	"testing"
	"time"

	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func TestDockerBackend_CreateContainer(t *testing.T) {
	m, c := prepareManager(t)

	t.Run("存在しないイメージを指定", func(t *testing.T) {
		t.Parallel()
		app := domain.Application{
			ID:        "test",
			UpdatedAt: time.Now(),
			Config: domain.ApplicationConfig{
				BuildConfig: &domain.BuildConfigRuntimeBuildpack{},
			},
		}
		st := domain.DesiredState{
			Runtime: []*domain.RuntimeDesiredState{{
				App:       &app,
				ImageName: "not-found",
				ImageTag:  "latest",
			}},
		}
		err := m.Synchronize(context.Background(), &st)
		assert.NoError(t, err) // fail-safe
	})

	t.Run("コンテナを正常に作成して起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"

		app := domain.Application{
			ID:        appID,
			UpdatedAt: time.Now(),
			Config: domain.ApplicationConfig{
				BuildConfig: &domain.BuildConfigRuntimeBuildpack{},
			},
		}
		st := domain.DesiredState{
			Runtime: []*domain.RuntimeDesiredState{{
				App:       &app,
				ImageName: image,
				ImageTag:  "latest",
			}},
		}
		err := m.Synchronize(context.Background(), &st)
		require.NoError(t, err)

		res, err := c.ContainerInspect(context.Background(), containerName(appID), client.ContainerInspectOptions{})
		require.NoError(t, err)

		assert.Equal(t, res.Container.Config.Image, image+":latest")
		assert.Equal(t, res.Container.Config.Labels[appLabel], "true")
		assert.Equal(t, res.Container.Config.Labels[appIDLabel], appID)

		_, err = c.ContainerRemove(context.Background(), res.Container.ID, client.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		require.NoError(t, err)
	})

	t.Run("コンテナを正常に作成して起動 (Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pij0bij90j20"

		app := domain.Application{
			ID:        appID,
			UpdatedAt: time.Now(),
			Config: domain.ApplicationConfig{
				BuildConfig: &domain.BuildConfigRuntimeBuildpack{},
			},
		}
		st := domain.DesiredState{
			Runtime: []*domain.RuntimeDesiredState{{
				App:       &app,
				ImageName: image,
				ImageTag:  "latest",
			}},
		}
		err := m.Synchronize(context.Background(), &st)
		require.NoError(t, err)

		app.UpdatedAt = time.Now() // Restart
		err = m.Synchronize(context.Background(), &st)
		require.NoError(t, err)

		res, err := c.ContainerInspect(context.Background(), containerName(appID), client.ContainerInspectOptions{})
		require.NoError(t, err)

		assert.Equal(t, res.Container.Config.Image, image+":latest")
		assert.Equal(t, res.Container.Config.Labels[appLabel], "true")
		assert.Equal(t, res.Container.Config.Labels[appIDLabel], appID)

		_, err = c.ContainerRemove(context.Background(), res.Container.ID, client.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		require.NoError(t, err)
	})
}
