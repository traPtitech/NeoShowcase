package k8simpl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

func (m *Manager) Restart(ctx context.Context, args container.RestartArgs) (*container.RestartResult, error) {
	data, _ := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				deploymentRestartAnnotation: time.Now().Format(time.RFC3339),
			},
		},
	})
	_, err := m.clientset.CoreV1().Pods(appNamespace).Patch(ctx, deploymentName(args.ApplicationID, args.EnvironmentID), types.MergePatchType, data, metav1.PatchOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to restart pod: %w", err)
	}
	return &container.RestartResult{}, nil
}
