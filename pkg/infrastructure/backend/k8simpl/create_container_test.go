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
	m, c := prepareManager(t, eventbus.NewLocal(hub.New()))

	t.Run("Podを正常に起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
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

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "test.localhost",
				Port:   80,
			},
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP, Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "98ygtfjfjhgj"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "ji9876fgoh.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
		}

		err = m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "bbbbbb.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID), metav1.DeleteOptions{}))
		}
	})
}
