package k8simpl

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) unregisterIngress(ctx context.Context, _ *domain.Application, website *domain.Website) error {
	name := serviceName(website.FQDN)
	err := b.clientset.NetworkingV1().Ingresses(appNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete ingress: %w", err)
	}
	return nil
}
