package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) unregisterService(ctx context.Context, _ *domain.Application, website *domain.Website) error {
	name := serviceName(website)
	err := b.client.CoreV1().Services(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

func (b *k8sBackend) unregisterMiddleware(ctx context.Context, _ *domain.Application, website *domain.Website) error {
	name := stripMiddlewareName(website)
	err := b.traefikClient.Middlewares(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete middleware: %w", err)
	}
	return nil
}

func (b *k8sBackend) unregisterIngress(ctx context.Context, app *domain.Application, website *domain.Website) error {
	name := serviceName(website)
	err := b.traefikClient.IngressRoutes(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete IngressRoute: %w", err)
	}
	if website.PathPrefix != "/" {
		err = b.unregisterMiddleware(ctx, app, website)
		if err != nil {
			return err
		}
	}
	return nil
}
