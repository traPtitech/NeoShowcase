package k8simpl

import (
	"fmt"
	"time"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	"github.com/traefik/traefik/v2/pkg/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

func (b *k8sBackend) runtimeService(app *domain.Application, website *domain.Website) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website),
			Namespace: b.config.Namespace,
			Labels:    b.appLabel(app.ID),
		},
		Spec: v1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: appSelector(app.ID),
			Ports: []v1.ServicePort{{
				Protocol:   "TCP",
				Port:       80,
				TargetPort: intstr.FromInt(website.HTTPPort),
			}},
		},
	}
}

func (b *k8sBackend) runtimeServiceRef(_ *domain.Application, website *domain.Website) []traefikv1alpha1.Service {
	return []traefikv1alpha1.Service{{
		LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
			Name:      serviceName(website),
			Kind:      "Service",
			Namespace: b.config.Namespace,
			Port:      intstr.FromInt(80),
			Scheme:    "http",
		},
	}}
}

func (b *k8sBackend) stripMiddleware(app *domain.Application, website *domain.Website) *traefikv1alpha1.Middleware {
	return &traefikv1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Middleware",
			APIVersion: "traefik.containo.us/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      stripMiddlewareName(website),
			Namespace: b.config.Namespace,
			Labels:    b.appLabel(app.ID),
		},
		Spec: traefikv1alpha1.MiddlewareSpec{
			StripPrefix: &dynamic.StripPrefix{
				Prefixes: []string{website.PathPrefix},
			},
		},
	}
}

func (b *k8sBackend) certificate(targetDomain string) *certmanagerv1.Certificate {
	return &certmanagerv1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      certificateName(targetDomain),
			Namespace: b.config.Namespace,
			Labels:    b.generalLabel(), // certificate may be shared by one or more apps
		},
		Spec: certmanagerv1.CertificateSpec{
			SecretName:  tlsSecretName(targetDomain),
			Duration:    &metav1.Duration{Duration: 90 * 24 * time.Hour},
			RenewBefore: &metav1.Duration{Duration: 15 * 24 * time.Hour},
			DNSNames:    []string{targetDomain},
			IssuerRef: certmetav1.ObjectReference{
				Name: b.config.TLS.CertManager.Issuer.Name,
				Kind: b.config.TLS.CertManager.Issuer.Kind,
			},
		},
	}
}

func (b *k8sBackend) ingressRoute(
	app *domain.Application,
	website *domain.Website,
	serviceRefs []traefikv1alpha1.Service,
	ads domain.AvailableDomainSlice,
) (
	*traefikv1alpha1.IngressRoute,
	[]*traefikv1alpha1.Middleware,
	[]*certmanagerv1.Certificate,
) {
	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, web.TraefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, web.TraefikHTTPEntrypoint)
	}

	var middlewareRefs []traefikv1alpha1.MiddlewareRef
	switch website.Authentication {
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
			middleware := b.stripMiddleware(app, website)
			middlewares = append(middlewares, middleware)
			middlewareRefs = append(middlewareRefs, traefikv1alpha1.MiddlewareRef{Name: middleware.Name})
		}
	}

	var tls *traefikv1alpha1.TLS
	var certs []*certmanagerv1.Certificate
	if website.HTTPS {
		if b.config.TLS.Type == tlsTypeTraefik {
			targetDomain := domain.TLSTargetDomain(b.config.TLS.Traefik.Wildcard, website, ads)
			tls = &traefikv1alpha1.TLS{
				SecretName:   tlsSecretName(targetDomain),
				CertResolver: b.config.TLS.Traefik.CertResolver,
				Domains:      []types.Domain{{Main: targetDomain}},
			}
		} else if b.config.TLS.Type == tlsTypeCertManager {
			targetDomain := domain.TLSTargetDomain(b.config.TLS.CertManager.Wildcard, website, ads)
			tls = &traefikv1alpha1.TLS{
				SecretName: tlsSecretName(targetDomain),
				Domains:    []types.Domain{{Main: targetDomain}},
			}
			certs = append(certs, b.certificate(targetDomain))
		}
	}

	ingressRoute := &traefikv1alpha1.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressRoute",
			APIVersion: "traefik.containo.us/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(website),
			Namespace: b.config.Namespace,
			Labels:    b.appLabel(app.ID),
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

	return ingressRoute, middlewares, certs
}
