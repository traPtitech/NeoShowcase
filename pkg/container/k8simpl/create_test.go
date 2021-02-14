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
}
