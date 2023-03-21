package k8simpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (b *k8sBackend) RestartContainer(ctx context.Context, appID string) error {
	data, _ := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				deploymentRestartAnnotation: time.Now().Format(time.RFC3339),
			},
		},
	})
	_, err := b.client.CoreV1().Pods(appNamespace).Patch(ctx, deploymentName(appID), types.MergePatchType, data, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to restart pod: %w", err)
	}
	return nil
}
