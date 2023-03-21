package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) unregisterService(ctx context.Context, _ *domain.Application, website *domain.Website) error {
	name := serviceName(website.FQDN)
	err := b.client.CoreV1().Services(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

func (b *k8sBackend) unregisterIngress(ctx context.Context, _ *domain.Application, website *domain.Website) error {
	name := serviceName(website.FQDN)
	err := b.traefikClient.IngressRoutes(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete IngressRoute: %w", err)
	}
	return nil
}
