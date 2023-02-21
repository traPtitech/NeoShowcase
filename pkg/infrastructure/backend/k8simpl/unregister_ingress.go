package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (b *k8sBackend) UnregisterIngress(ctx context.Context, appID string) error {
	name := deploymentName(appID)
	err := b.clientset.NetworkingV1().Ingresses(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete ingress: %w", err)
	}
	return nil
}
