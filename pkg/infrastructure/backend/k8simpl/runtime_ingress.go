package k8simpl

import (
	"fmt"

	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	"github.com/traefik/traefik/v2/pkg/types"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

func runtimeService(app *domain.Application, website *domain.Website) *apiv1.Service {
	return &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website),
			Namespace: appNamespace,
			Labels:    resourceLabels(app.ID),
		},
		Spec: apiv1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: resourceLabels(app.ID),
			Ports: []apiv1.ServicePort{{
				Protocol:   "TCP",
				Port:       80,
				TargetPort: intstr.FromInt(website.HTTPPort),
			}},
		},
	}
}

func runtimeServiceRef(_ *domain.Application, website *domain.Website) []traefikv1alpha1.Service {
	return []traefikv1alpha1.Service{{
		LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
			Name:      serviceName(website),
			Kind:      "Service",
			Namespace: appNamespace,
			Port:      intstr.FromInt(80),
			Scheme:    "http",
		},
	}}
}

func stripMiddleware(_ *domain.Application, website *domain.Website, labels map[string]string) *traefikv1alpha1.Middleware {
	return &traefikv1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Middleware",
			APIVersion: "traefik.containo.us/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      stripMiddlewareName(website),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: traefikv1alpha1.MiddlewareSpec{
			StripPrefix: &dynamic.StripPrefix{
				Prefixes: []string{website.PathPrefix},
			},
		},
	}
}

func ingressRoute(app *domain.Application, website *domain.Website, labels map[string]string, serviceRefs []traefikv1alpha1.Service) (*traefikv1alpha1.IngressRoute, []*traefikv1alpha1.Middleware) {
	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, web.TraefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, web.TraefikHTTPEntrypoint)
	}

	var middlewareRefs []traefikv1alpha1.MiddlewareRef
	switch app.Config.Authentication {
	case domain.AuthenticationTypeSoft:
		middlewareRefs = append(middlewareRefs,
			traefikv1alpha1.MiddlewareRef{Name: web.TraefikAuthSoftMiddleware},
			traefikv1alpha1.MiddlewareRef{Name: web.TraefikAuthMiddleware},
		)
	case domain.AuthenticationTypeHard:
		middlewareRefs = append(middlewareRefs,
			traefikv1alpha1.MiddlewareRef{Name: web.TraefikAuthHardMiddleware},
			traefikv1alpha1.MiddlewareRef{Name: web.TraefikAuthMiddleware},
		)
	}

	var rule string
	var middlewares []*traefikv1alpha1.Middleware
	if website.PathPrefix == "/" {
		rule = fmt.Sprintf("Host(`%s`)", website.FQDN)
	} else {
		rule = fmt.Sprintf("Host(`%s`) && PathPrefix(`%s`)", website.FQDN, website.PathPrefix)
		if website.StripPrefix {
			middleware := stripMiddleware(app, website, labels)
			middlewares = append(middlewares, middleware)
			middlewareRefs = append(middlewareRefs, traefikv1alpha1.MiddlewareRef{Name: middleware.Name})
		}
	}

	var tls *traefikv1alpha1.TLS
	if website.HTTPS {
		tls = &traefikv1alpha1.TLS{
			SecretName:   tlsSecretName(website.FQDN),
			CertResolver: web.TraefikCertResolver,
			Domains:      []types.Domain{{Main: website.FQDN}},
		}
	}

	ingressRoute := &traefikv1alpha1.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressRoute",
			APIVersion: "traefik.containo.us/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: traefikv1alpha1.IngressRouteSpec{
			EntryPoints: entrypoints,
			Routes: []traefikv1alpha1.Route{{
				Match:       rule,
				Kind:        "Rule",
				Services:    serviceRefs,
				Middlewares: middlewareRefs,
			}},
			TLS: tls,
		},
	}

	return ingressRoute, middlewares
}
