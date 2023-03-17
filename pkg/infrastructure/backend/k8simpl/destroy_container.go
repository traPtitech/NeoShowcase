package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	err := b.clientset.CoreV1().Pods(appNamespace).Delete(ctx, deploymentName(app.ID), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete pod: %w", err)
	}

	for _, website := range app.Websites {
		err = b.clientset.CoreV1().Services(appNamespace).Delete(ctx, serviceName(website.FQDN), metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("failed to delete service: %w", err)
		}
		err = b.unregisterIngress(ctx, app, website)
		if err != nil {
			return err
		}
	}

	return nil
}
