package k8simpl

import (
	"context"
	"fmt"

	"github.com/volatiletech/null/v8"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (b *k8sBackend) RegisterIngress(ctx context.Context, appID string, envID string, host string, destination null.String, port null.Int) error {
	labels := map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: appID,
		appContainerEnvironmentIDLabel: envID,
	}

	svc := &networkingv1.IngressServiceBackend{
		Name: deploymentName(appID, envID),
		Port: networkingv1.ServiceBackendPort{
			Number: 80,
		},
	}
	if destination.Valid {
		svc.Name = destination.String
	}
	if port.Valid {
		svc.Port.Number = int32(port.Int)
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(appID, envID),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: pathTypePtr(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: svc,
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
