package k8simpl

import (
	"context"
	"testing"

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func TestK8sBackend_CreateContainer(t *testing.T) {
	m, c, tc := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("Podを正常に起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"

		app := domain.Application{
			ID: appID,
		}
		err := m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
			ImageName: image,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "pijojopjnnna"

		website := &domain.Website{
			FQDN:       "test.localhost",
			PathPrefix: "/",
			HTTPPort:   80,
		}
		app := domain.Application{
			ID:       appID,
			Websites: []*domain.Website{website},
		}
		err := m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
			ImageName: image,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), serviceName(website), metav1.DeleteOptions{}))
			require.NoError(t, tc.IngressRoutes(appNamespace).Delete(context.Background(), serviceName(website), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP, Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "98ygtfjfjhgj"

		website := &domain.Website{
			FQDN:       "ji9876fgoh.localhost",
			PathPrefix: "/test",
			HTTPPort:   80,
		}
		app := domain.Application{
			ID:       appID,
			Websites: []*domain.Website{website},
		}
		err := m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
			ImageName: image,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
		}

		err = m.CreateContainer(context.Background(), &app, domain.ContainerCreateArgs{
			ImageName: image,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), serviceName(website), metav1.DeleteOptions{}))
			require.NoError(t, tc.IngressRoutes(appNamespace).Delete(context.Background(), serviceName(website), metav1.DeleteOptions{}))
			require.NoError(t, tc.Middlewares(appNamespace).Delete(context.Background(), stripMiddlewareName(website), metav1.DeleteOptions{}))
		}
	})
}
