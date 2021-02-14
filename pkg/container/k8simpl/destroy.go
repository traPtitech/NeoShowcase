package k8simpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) Destroy(ctx context.Context, args container.DestroyArgs) (*container.DestroyResult, error) {
	err := m.clientset.CoreV1().Pods(appNamespace).Delete(ctx, deploymentName(args.ApplicationID, args.EnvironmentID), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to delete pod: %w", err)
	}

	err = m.clientset.CoreV1().Services(appNamespace).Delete(ctx, deploymentName(args.ApplicationID, args.EnvironmentID), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to delete service: %w", err)
	}

	err = m.clientset.NetworkingV1beta1().Ingresses(appNamespace).Delete(ctx, deploymentName(args.ApplicationID, args.EnvironmentID), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to delete ingress: %w", err)
	}

	return &container.DestroyResult{}, nil
}
