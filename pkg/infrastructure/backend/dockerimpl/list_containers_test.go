package dockerimpl

import (
	"context"
	"strconv"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func TestDockerBackend_ListContainers(t *testing.T) {
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("空リスト", func(t *testing.T) {
		result, err := m.ListContainers(context.Background())
		if assert.NoError(t, err) {
			assert.Empty(t, result)
		}
	})

	t.Run("正常", func(t *testing.T) {
		image := "hello-world"
		baseAppID := "afiowjiodncx"
		n := 5
		for i := 0; i < n; i++ {
			i := i
			app := domain.Application{
				ID: baseAppID + strconv.Itoa(i),
			}
			err := m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
				ImageName: image,
			})
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = c.RemoveContainer(docker.RemoveContainerOptions{
					ID:            containerName(app.ID),
					RemoveVolumes: true,
					Force:         true,
				})
			})
		}

		result, err := m.ListContainers(context.Background())
		if assert.NoError(t, err) {
			assert.Len(t, result, n)
		}
	})
}
