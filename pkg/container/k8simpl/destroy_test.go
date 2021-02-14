package k8simpl

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"testing"
)

func TestManager_Destroy(t *testing.T) {
	m, c, _ := prepareManager(t)

	t.Run("Podを正常に削除", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "ap8ajievpjap"
		envID := "24j0jadlskfj"

		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		require.NoError(t, err)
		waitPodRunning(t, c, deploymentName(appID, envID))

		_, err = m.Destroy(context.Background(), container.DestroyArgs{
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		assert.NoError(t, err)
	})

	t.Run("Podを正常に削除 (HTTP)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "pjpoi2efeioji"
		envID := "908uyoinkasdf"

		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &container.HTTPProxy{
				Domain: "test.localhost",
				Port:   80,
			},
		})
		require.NoError(t, err)
		waitPodRunning(t, c, deploymentName(appID, envID))

		_, err = m.Destroy(context.Background(), container.DestroyArgs{
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		assert.NoError(t, err)
	})
}
