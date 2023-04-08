package k8simpl

import (
	"context"
	"testing"
	"time"

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/require"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func TestK8sBackend_Synchronize(t *testing.T) {
	const appNamespace = "neoshowcase-apps"

	m, c, tc := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("Podを正常に起動", func(t *testing.T) {
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"

		app := domain.Application{
			ID:        appID,
			UpdatedAt: time.Now(),
		}
		err := m.SynchronizeRuntime(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
		}})
		require.NoError(t, err)
		exists[*appsv1.StatefulSet](t, deploymentName(appID), c.AppsV1().StatefulSets(appNamespace))
		waitPodRunning(t, m, appID)

		err = m.SynchronizeRuntime(context.Background(), nil)
		require.NoError(t, err)
		waitPodDeleted(t, m, appID) // NOTE: foreground delete
		notExists[*appsv1.StatefulSet](t, deploymentName(appID), c.AppsV1().StatefulSets(appNamespace))
	})

	t.Run("Podを正常に起動 (HTTP)", func(t *testing.T) {
		t.SkipNow()
		image := "chussenot/tiny-server"
		appID := "pijojopjnnna"

		website := &domain.Website{
			ID:          "282d4394a71686dcc4a3e2",
			FQDN:        "test.localhost",
			PathPrefix:  "/test",
			StripPrefix: false,
			HTTPPort:    80,
		}
		app := domain.Application{
			ID:        appID,
			Websites:  []*domain.Website{website},
			UpdatedAt: time.Now(),
		}
		err := m.SynchronizeRuntime(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
		}})
		require.NoError(t, err)
		exists[*appsv1.StatefulSet](t, deploymentName(appID), c.AppsV1().StatefulSets(appNamespace))
		exists[*corev1.Service](t, serviceName(website), c.CoreV1().Services(appNamespace))
		exists[*traefikv1alpha1.IngressRoute](t, serviceName(website), tc.IngressRoutes(appNamespace))
		waitPodRunning(t, m, appID)

		err = m.SynchronizeRuntime(context.Background(), nil)
		require.NoError(t, err)
		waitPodDeleted(t, m, appID) // NOTE: foreground delete
		notExists[*appsv1.StatefulSet](t, deploymentName(appID), c.AppsV1().StatefulSets(appNamespace))
		notExists[*corev1.Service](t, serviceName(website), c.CoreV1().Services(appNamespace))
		notExists[*traefikv1alpha1.IngressRoute](t, serviceName(website), tc.IngressRoutes(appNamespace))
	})

	t.Run("Podを正常に起動 (HTTP, Recreate)", func(t *testing.T) {
		t.SkipNow()
		image := "chussenot/tiny-server"
		appID := "98ygtfjfjhgj"

		website := &domain.Website{
			ID:          "a3fd3e4df5d66bfcb8f11c",
			FQDN:        "ji9876fgoh.localhost",
			PathPrefix:  "/test",
			StripPrefix: true,
			HTTPPort:    80,
		}
		app := domain.Application{
			ID:        appID,
			Websites:  []*domain.Website{website},
			UpdatedAt: time.Now(),
		}
		err := m.SynchronizeRuntime(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
		}})
		require.NoError(t, err)
		waitPodRunning(t, m, appID)

		app.UpdatedAt = time.Now() // Restart
		err = m.SynchronizeRuntime(context.Background(), []*domain.AppDesiredState{{
			App:       &app,
			ImageName: image,
			ImageTag:  "latest",
		}})
		require.NoError(t, err)
		exists[*appsv1.StatefulSet](t, deploymentName(appID), c.AppsV1().StatefulSets(appNamespace))
		exists[*corev1.Service](t, serviceName(website), c.CoreV1().Services(appNamespace))
		exists[*traefikv1alpha1.IngressRoute](t, serviceName(website), tc.IngressRoutes(appNamespace))
		exists[*traefikv1alpha1.Middleware](t, stripMiddlewareName(website), tc.Middlewares(appNamespace))
		waitPodRunning(t, m, appID)

		err = m.SynchronizeRuntime(context.Background(), nil)
		require.NoError(t, err)
		waitPodDeleted(t, m, appID) // NOTE: foreground delete
		notExists[*appsv1.StatefulSet](t, deploymentName(appID), c.AppsV1().StatefulSets(appNamespace))
		notExists[*corev1.Service](t, serviceName(website), c.CoreV1().Services(appNamespace))
		notExists[*traefikv1alpha1.IngressRoute](t, serviceName(website), tc.IngressRoutes(appNamespace))
		notExists[*traefikv1alpha1.Middleware](t, stripMiddlewareName(website), tc.Middlewares(appNamespace))
	})
}
