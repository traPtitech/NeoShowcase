package k8simpl

import (
	"context"
	"fmt"

	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"github.com/traefik/traefik/v2/pkg/types"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) registerService(ctx context.Context, _ *domain.Application, website *domain.Website, podSelector map[string]string) error {
	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website.FQDN),
			Namespace: appNamespace,
			Labels:    map[string]string{},
		},
		Spec: apiv1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: podSelector,
			Ports: []apiv1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       80,
					TargetPort: intstr.FromInt(website.HTTPPort),
				},
			},
		},
	}

	if _, err := b.client.CoreV1().Services(appNamespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			if err = b.client.CoreV1().Services(appNamespace).Delete(ctx, svc.Name, metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("failed to delete service: %w", err)
			}
			if _, err := b.client.CoreV1().Services(appNamespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create service: %w", err)
		}
	}

	return nil
}

func (b *k8sBackend) registerIngress(ctx context.Context, app *domain.Application, website *domain.Website) error {
	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, traefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, traefikHTTPEntrypoint)
	}

	var middlewares []v1alpha1.MiddlewareRef
	switch app.Config.Authentication {
	case domain.AuthenticationTypeSoft:
		middlewares = append(middlewares,
			v1alpha1.MiddlewareRef{Name: traefikAuthSoftMiddleware},
			v1alpha1.MiddlewareRef{Name: traefikAuthMiddleware},
		)
	case domain.AuthenticationTypeHard:
		middlewares = append(middlewares,
			v1alpha1.MiddlewareRef{Name: traefikAuthHardMiddleware},
			v1alpha1.MiddlewareRef{Name: traefikAuthMiddleware},
		)
	}

	var tls *v1alpha1.TLS
	if website.HTTPS {
		tls = &v1alpha1.TLS{
			SecretName:   serviceName(website.FQDN),
			CertResolver: traefikCertResolver,
			Domains:      []types.Domain{{Main: website.FQDN}},
		}
	}

	ingressRoute := &v1alpha1.IngressRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website.FQDN),
			Namespace: appNamespace,
			Labels:    map[string]string{},
		},
		Spec: v1alpha1.IngressRouteSpec{
			EntryPoints: entrypoints,
			Routes: []v1alpha1.Route{{
				Match:       fmt.Sprintf("Host(`%s`)", website.FQDN),
				Kind:        "Rule",
				Middlewares: middlewares,
				Services: []v1alpha1.Service{{
					LoadBalancerSpec: v1alpha1.LoadBalancerSpec{
						Name:      serviceName(website.FQDN),
						Kind:      "Service",
						Namespace: appNamespace,
						Port:      intstr.FromInt(website.HTTPPort),
						Scheme:    "http",
					},
				}},
			}},
			TLS: tls,
		},
	}

	if _, err := b.traefikClient.IngressRoutes(appNamespace).Create(ctx, ingressRoute, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			if err = b.traefikClient.IngressRoutes(appNamespace).Delete(ctx, ingressRoute.Name, metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("failed to delete IngressRoute: %w", err)
			}
			if _, err = b.traefikClient.IngressRoutes(appNamespace).Create(ctx, ingressRoute, metav1.CreateOptions{}); err != nil {
				return fmt.Errorf("failed to create IngressRoute: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create IngressRoute: %w", err)
		}
	}
	return nil
}
