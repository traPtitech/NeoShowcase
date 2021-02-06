package dockerimpl

import (
	"context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"testing"
)

func TestManager_Destroy(t *testing.T) {
	m, c, _ := prepareManager(t)

	t.Run("存在しないコンテナを指定", func(t *testing.T) {
		t.Parallel()
		_, err := m.Destroy(context.Background(), container.DestroyArgs{
			ApplicationID: "notfound",
			EnvironmentID: "notfound",
		})
		assert.Error(t, err)
	})

	t.Run("コンテナを正常に削除", func(t *testing.T) {
		t.Parallel()
		appID := "ojoionaonidp"
		envID := "bhhfkadajlkh"
		_, err := c.CreateContainer(docker.CreateContainerOptions{
			Name: containerName(appID, envID),
			Config: &docker.Config{
				Image: "alpine:latest",
			},
		})
		require.NoError(t, err)

		_, err = m.Destroy(context.Background(), container.DestroyArgs{
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		assert.NoError(t, err)
	})

	t.Run("稼働中のコンテナを削除", func(t *testing.T) {
		t.Parallel()
		appID := "pjipjjijoinn"
		envID := "wefadsnaiomo"
		cont, err := c.CreateContainer(docker.CreateContainerOptions{
			Name: containerName(appID, envID),
			Config: &docker.Config{
				Image: "alpine:latest",
				Cmd:   []string{"sleep", "100"},
			},
		})
		require.NoError(t, err)
		require.NoError(t, c.StartContainer(cont.ID, nil))

		_, err = m.Destroy(context.Background(), container.DestroyArgs{
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		assert.NoError(t, err)
	})
}
