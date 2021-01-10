package k8simpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	autoscalev1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) Stop(ctx context.Context, args container.StopArgs) (*container.StopResult, error) {
	_, err := m.clientset.AppsV1().Deployments(appNamespace).UpdateScale(ctx, deploymentName(args.ApplicationID), &autoscalev1.Scale{Spec: autoscalev1.ScaleSpec{Replicas: 0}}, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to scale deployment to 0: %w", err)
	}
	return &container.StopResult{}, nil
}
