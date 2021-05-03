package k8simpl

import (
	"context"
	"testing"

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func TestK8sBackend_DestroyContainer(t *testing.T) {
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("Podを正常に削除", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "ap8ajievpjap"
		envID := "24j0jadlskfj"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		require.NoError(t, err)
		waitPodRunning(t, c, deploymentName(appID, envID))

		err = m.DestroyContainer(context.Background(), appID, envID)
		assert.NoError(t, err)
	})

	t.Run("Podを正常に削除 (HTTP)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "pjpoi2efeioji"
		envID := "908uyoinkasdf"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "test.localhost",
				Port:   80,
			},
		})
		require.NoError(t, err)
		waitPodRunning(t, c, deploymentName(appID, envID))

		err = m.DestroyContainer(context.Background(), appID, envID)
		assert.NoError(t, err)
	})
}
