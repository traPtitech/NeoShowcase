package k8simpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	autoscalev1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) Start(ctx context.Context, args container.StartArgs) (*container.StartResult, error) {
	_, err := m.clientset.AppsV1().Deployments(appNamespace).UpdateScale(ctx, deploymentName(args.ApplicationID, args.EnvironmentID), &autoscalev1.Scale{Spec: autoscalev1.ScaleSpec{Replicas: 1}}, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to scale deployment to 1: %w", err)
	}
	return &container.StartResult{}, nil
}
