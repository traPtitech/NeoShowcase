package k8simpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) Destroy(ctx context.Context, args container.DestroyArgs) (*container.DestroyResult, error) {
	err := m.clientset.AppsV1().Deployments(appNamespace).Delete(ctx, deploymentName(args.ApplicationID), metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete deployment: %w", err)
	}

	err = m.clientset.CoreV1().Services(appNamespace).Delete(ctx, deploymentName(args.ApplicationID), metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete service: %w", err)
	}

	err = m.clientset.NetworkingV1beta1().Ingresses(appNamespace).Delete(ctx, deploymentName(args.ApplicationID), metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete ingress: %w", err)
	}

	return &container.DestroyResult{}, nil
}
