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
	m, _, _ := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("Podを正常に削除", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "ap8ajievpjap"

		app := domain.Application{
			ID: appID,
		}
		err := m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
			ImageName: image,
		})
		require.NoError(t, err)
		waitPodRunning(t, m, appID)

		err = m.DestroyContainer(context.Background(), &app)
		assert.NoError(t, err)
	})

	t.Run("Podを正常に削除 (HTTP)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "pjpoi2efeioji"

		app := domain.Application{
			ID: appID,
			Websites: []*domain.Website{{
				FQDN:       "test.localhost",
				PathPrefix: "/",
				HTTPPort:   80,
			}},
		}
		err := m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
			ImageName: image,
		})
		require.NoError(t, err)
		waitPodRunning(t, m, appID)

		err = m.DestroyContainer(context.Background(), &app)
		assert.NoError(t, err)
	})
}
