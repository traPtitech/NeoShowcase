package k8simpl

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestManager_Create(t *testing.T) {
	m, c, _ := prepareManager(t)

	t.Run("Podを正常に起動", func(t *testing.T) {
		t.Parallel()
		image := "tianon/sleeping-beauty"
		appID := "pjpjpjoijion"
		envID := "fewwfadsface"

		_, err := m.Create(context.Background(), container.CreateArgs{
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

		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &container.HTTPProxy{
				Domain: "test.localhost",
				Port:   80,
			},
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.NetworkingV1beta1().Ingresses(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP, Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "98ygtfjfjhgj"
		envID := "wertyuyui987"

		_, err := m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &container.HTTPProxy{
				Domain: "ji9876fgoh.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
		}

		_, err = m.Create(context.Background(), container.CreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			EnvironmentID: envID,
			HTTPProxy: &container.HTTPProxy{
				Domain: "bbbbbb.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, envID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
			require.NoError(t, c.NetworkingV1beta1().Ingresses(appNamespace).Delete(context.Background(), deploymentName(appID, envID), metav1.DeleteOptions{}))
		}
	})

}
