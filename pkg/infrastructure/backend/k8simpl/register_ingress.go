package k8simpl

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) registerIngress(ctx context.Context, app *domain.Application, website *domain.Website) error {
	labels := map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: app.ID,
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website.FQDN),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: website.FQDN,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: pathTypePtr(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName(website.FQDN),
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if _, err := b.clientset.NetworkingV1().Ingresses(appNamespace).Create(ctx, ingress, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			if _, err = b.clientset.NetworkingV1().Ingresses(appNamespace).Update(ctx, ingress, metav1.UpdateOptions{}); err != nil {
				return fmt.Errorf("failed to update ingress: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create ingress: %w", err)
		}
	}
	return nil
}
