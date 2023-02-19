package dockerimpl

import (
	"context"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func TestDockerBackend_DestroyContainer(t *testing.T) {
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("存在しないコンテナを指定", func(t *testing.T) {
		t.Parallel()
		err := m.DestroyContainer(context.Background(), "notfound")
		assert.Error(t, err)
	})

	t.Run("コンテナを正常に削除", func(t *testing.T) {
		t.Parallel()
		appID := "ojoionaonidp"
		_, err := c.CreateContainer(docker.CreateContainerOptions{
			Name: containerName(appID),
			Config: &docker.Config{
				Image: "alpine:latest",
			},
		})
		require.NoError(t, err)

		err = m.DestroyContainer(context.Background(), appID)
		assert.NoError(t, err)
	})

	t.Run("稼働中のコンテナを削除", func(t *testing.T) {
		t.Parallel()
		appID := "pjipjjijoinn"
		cont, err := c.CreateContainer(docker.CreateContainerOptions{
			Name: containerName(appID),
			Config: &docker.Config{
				Image: "alpine:latest",
				Cmd:   []string{"sleep", "100"},
			},
		})
		require.NoError(t, err)
		require.NoError(t, c.StartContainer(cont.ID, nil))

		err = m.DestroyContainer(context.Background(), appID)
		assert.NoError(t, err)
	})
}
