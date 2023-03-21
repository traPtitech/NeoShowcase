package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	for _, website := range app.Websites {
		err := b.unregisterIngress(ctx, app, website)
		if err != nil {
			return err
		}
		err = b.unregisterService(ctx, app, website)
		if err != nil {
			return err
		}
	}

	err := b.client.CoreV1().Pods(appNamespace).Delete(ctx, deploymentName(app.ID), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete pod: %w", err)
	}
	return nil
}
