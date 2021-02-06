package dockerimpl

import (
	"context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"strconv"
	"testing"
)

func TestManager_List(t *testing.T) {
	m, c, _ := prepareManager(t)

	t.Run("空リスト", func(t *testing.T) {
		result, err := m.List(context.Background())
		if assert.NoError(t, err) {
			assert.Empty(t, result.Containers)
		}
	})

	t.Run("正常", func(t *testing.T) {
		image := "hello-world"
		appID := "afiowjiodncx"
		envID := "adhihpillomo"
		n := 5
		for i := 0; i < n; i++ {
			i := i
			_, err := m.Create(context.Background(), container.CreateArgs{
				ImageName:     image,
				ApplicationID: appID + strconv.Itoa(i),
				EnvironmentID: envID + strconv.Itoa(i),
			})
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = c.RemoveContainer(docker.RemoveContainerOptions{
					ID:            containerName(appID+strconv.Itoa(i), envID+strconv.Itoa(i)),
					RemoveVolumes: true,
					Force:         true,
				})
			})
		}

		result, err := m.List(context.Background())
		if assert.NoError(t, err) {
			assert.Len(t, result.Containers, n)
		}
	})
}
