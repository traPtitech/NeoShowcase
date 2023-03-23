package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	err := b.destroyRuntimeIngresses(ctx, app)
	if err != nil {
		return fmt.Errorf("failed to destroy runtime ingress resources: %w", err)
	}

	err = b.client.AppsV1().StatefulSets(appNamespace).Delete(ctx, deploymentName(app.ID), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete pod: %w", err)
	}
	return nil
}
