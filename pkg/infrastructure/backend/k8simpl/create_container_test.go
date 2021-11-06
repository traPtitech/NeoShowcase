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
		branchID := "fewwfadsface"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			BranchID:      branchID,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, branchID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, branchID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "pijojopjnnna"
		branchID := "2io3isaoioij"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			BranchID:      branchID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "test.localhost",
				Port:   80,
			},
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, branchID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, branchID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID, branchID), metav1.DeleteOptions{}))
		}
	})

	t.Run("Podを正常に起動 (HTTP, Recreate)", func(t *testing.T) {
		t.Parallel()
		image := "chussenot/tiny-server"
		appID := "98ygtfjfjhgj"
		branchID := "wertyuyui987"

		err := m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			BranchID:      branchID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "ji9876fgoh.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, branchID))
		}

		err = m.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ImageName:     image,
			ApplicationID: appID,
			BranchID:      branchID,
			HTTPProxy: &domain.ContainerHTTPProxy{
				Domain: "bbbbbb.localhost",
				Port:   80,
			},
			Recreate: true,
		})
		if assert.NoError(t, err) {
			waitPodRunning(t, c, deploymentName(appID, branchID))
			require.NoError(t, c.CoreV1().Pods(appNamespace).Delete(context.Background(), deploymentName(appID, branchID), metav1.DeleteOptions{}))
			require.NoError(t, c.CoreV1().Services(appNamespace).Delete(context.Background(), deploymentName(appID, branchID), metav1.DeleteOptions{}))
		}
	})
}
