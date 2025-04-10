package dockerimpl

import (
	"context"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types/container"
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
		config := domain.ApplicationConfig{
			BuildConfig: &domain.BuildConfigRuntimeCmd{},
		}
		for i := range n {
			app := domain.Application{
				ID:     baseAppID + strconv.Itoa(i),
				Config: config,
			}
			apps = append(apps, &domain.RuntimeDesiredState{
				App:       &app,
				ImageName: image,
				ImageTag:  "latest",
			})
			t.Cleanup(func() {
				_ = c.ContainerRemove(context.Background(), containerName(app.ID), container.RemoveOptions{
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
