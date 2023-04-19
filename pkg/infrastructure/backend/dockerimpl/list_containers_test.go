package dockerimpl

import (
	"context"
	"strconv"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func TestDockerBackend_ListContainers(t *testing.T) {
	m, c := prepareManager(t)

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
		var apps []*domain.RuntimeDesiredState
		for i := 0; i < n; i++ {
			app := domain.Application{
				ID: baseAppID + strconv.Itoa(i),
			}
			apps = append(apps, &domain.RuntimeDesiredState{
				App:       &app,
				ImageName: image,
				ImageTag:  "latest",
			})
			t.Cleanup(func() {
				_ = c.RemoveContainer(docker.RemoveContainerOptions{
					ID:            containerName(app.ID),
					RemoveVolumes: true,
					Force:         true,
				})
			})
		}

		st := domain.DesiredState{Runtime: apps}
		err := m.Synchronize(context.Background(), &st)
		require.NoError(t, err)

		result, err := m.ListContainers(context.Background())
		if assert.NoError(t, err) {
			assert.Len(t, result, n)
		}
	})
}
