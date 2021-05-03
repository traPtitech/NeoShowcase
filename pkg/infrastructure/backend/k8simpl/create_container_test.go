package k8simpl

import (
	"context"
	"testing"

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestK8sBackend_CreateContainer(t *testing.T) {
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("Podを正常に起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"
		envID := "fewwfadsface"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "pijojopjnnna"
		envID := "2io3isaoioij"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "test.localhost",
				Port:   80,
			},
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.NetworkingV1().Ingresses(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP, Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "98ygtfjfjhgj"
		envID := "wertyuyui987"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "ji9876fgoh.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
		}

		err = m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "bbbbbb.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.NetworkingV1().Ingresses(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
		}
	})
}
