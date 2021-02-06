package dockerimpl

import (
	"context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"testing"
)

func TestManager_Create(t *testing.T) {
	m, c, _ := prepareManager(t)

	t.Run("存在しないイメージを指定", func(t *testing.T) {
		t.Parallel()
		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName: "notfoundimage",
		})
		assert.Error(t, err)
	})

	t.Run("コンテナを正常に作成", func(t *testing.T) {
		t.Parallel()
		image := "hello-world"
		appID := "afiowjiodncx"
		envID := "adhihpillomo"
		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			NoStart:       true,
		})
		if assert.NoError(t, err) {
			cont, err := c.InspectContainerWithOptions(docker.InspectContainerOptions{
				ID: containerName(appID, envID),
			})
			require.NoError(t, err)

			assert.Equal(t, cont.Config.Image, image+":latest")
			assert.Equal(t, cont.Config.Labels[appContainerLabel], "true")
			assert.Equal(t, cont.Config.Labels[appContainerApplicationIDLabel], appID)
			assert.Equal(t, cont.Config.Labels[appContainerEnvironmentIDLabel], envID)

			require.NoError(t, c.RemoveContainer(docker.RemoveContainerOptions{
				ID:            containerName(appID, envID),
				RemoveVolumes: true,
				Force:         true,
			}))
		}
	})

	t.Run("コンテナを正常に作成して起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"
		envID := "fewwfadsface"

		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		if assert.NoError(t, err) {
			cont, err := c.InspectContainerWithOptions(docker.InspectContainerOptions{
				ID: containerName(appID, envID),
			})
			require.NoError(t, err)

			assert.Equal(t, cont.Config.Image, image+":latest")
			assert.Equal(t, cont.Config.Labels[appContainerLabel], "true")
			assert.Equal(t, cont.Config.Labels[appContainerApplicationIDLabel], appID)
			assert.Equal(t, cont.Config.Labels[appContainerEnvironmentIDLabel], envID)

			require.NoError(t, c.RemoveContainer(docker.RemoveContainerOptions{
				ID:            containerName(appID, envID),
				RemoveVolumes: true,
				Force:         true,
			}))
		}
	})
}
