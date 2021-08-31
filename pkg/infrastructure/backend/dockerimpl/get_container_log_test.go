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

func TestDockerBackend_GetContainerLog(t *testing.T) {
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("正常", func(t *testing.T) {
		image := "hello-world"
		appID := "afiowjiodncx"
		envID := "adhihpillomo"
		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = c.RemoveContainer(docker.RemoveContainerOptions{
				ID:            containerName(appID, envID),
				RemoveVolumes: true,
				Force:         true,
			})
		})
		opts := LogOptions{}
		result, err := m.GetContainerStdOut(context.Background(), appID, envID, opts)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, result)
		}
	})
}
