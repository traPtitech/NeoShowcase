package k8simpl

import (
	"encoding/json"
	"fmt"
	"time"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	"github.com/traefik/traefik/v3/pkg/types"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

func (b *Backend) stripMiddleware(app *domain.Application, website *domain.Website) *traefikv1alpha1.Middleware {
	return &traefikv1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Middleware",
			APIVersion: "traefik.io/v1alpha1",
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

func (b *Backend) certificate(targetDomain string) *certmanagerv1.Certificate {
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
			SecretName: tlsSecretName(targetDomain),
			SecretTemplate: &certmanagerv1.CertificateSecretTemplate{
				Labels: b.generalLabelWithoutManagement(),
			},
			Duration:    &metav1.Duration{Duration: 90 * 24 * time.Hour},
			RenewBefore: &metav1.Duration{Duration: 30 * 24 * time.Hour},
			DNSNames:    []string{targetDomain},
			IssuerRef: certmetav1.ObjectReference{
				Name: b.config.TLS.CertManager.Issuer.Name,
				Kind: b.config.TLS.CertManager.Issuer.Kind,
			},
		},
	}
}

func (b *Backend) jsonSablierConfig(app *domain.Application) []byte {
	type DynamicConfig = struct {
		DisplayName string `json:"displayName"`
		ShowDetails string `json:"showDetails"`
		Theme       string `json:"theme"`
	}
	type BlockingConfig = struct {
		Timeout string `json:"timeout"`
	}
	type SablierConfig = struct {
		SablierURL      string          `json:"sablierURL"`
		Group           string          `json:"group"`
		SessionDuration string          `json:"sessionDuration"`
		Dynamic         *DynamicConfig  `json:"dynamic,omitempty"`
		Blocking        *BlockingConfig `json:"blocking,omitempty"`
	}

	config := SablierConfig{
		SablierURL:      b.config.Middleware.Sablier.SablierURL,
		Group:           sablierGroupName(app.ID),
		SessionDuration: b.config.Middleware.Sablier.SessionDuration,
	}

	switch app.Config.BuildConfig.GetRuntimeConfig().AutoShutdown.Startup {
	case domain.StartupBehaviorLoadingPage:
		config.Dynamic = &DynamicConfig{
			DisplayName: app.Name,
			ShowDetails: "true",
			Theme:       b.config.Middleware.Sablier.Dynamic.Theme,
		}
	case domain.StartupBehaviorBlocking:
		config.Blocking = &BlockingConfig{
			Timeout: b.config.Middleware.Sablier.Blocking.Timeout,
		}
	}

	data, _ := json.Marshal(config)
	return data
}

func (b *Backend) sablierMiddleware(app *domain.Application) *traefikv1alpha1.Middleware {
	configData := b.jsonSablierConfig(app)
	return &traefikv1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Middleware",
			APIVersion: "traefik.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      sablierMiddlewareName(app.ID),
			Namespace: b.config.Namespace,
			Labels:    b.appLabel(app.ID),
		},
		Spec: traefikv1alpha1.MiddlewareSpec{
			Plugin: map[string]v1.JSON{
				"sablier": {Raw: configData},
			},
		},
	}
}

func (b *Backend) ingressRoute(
	app *domain.Application,
	website *domain.Website,
	serviceRefs []traefikv1alpha1.Service,
) (
	*traefikv1alpha1.IngressRoute,
	[]*traefikv1alpha1.Middleware,
) {
	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, web.TraefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, web.TraefikHTTPEntrypoint)
	}

	var middlewareRefs []traefikv1alpha1.MiddlewareRef
	authConfig := b.targetAuth(website.FQDN)
	if authConfig != nil {
		switch website.Authentication {
		case domain.AuthenticationTypeSoft:
			for _, mw := range authConfig.Soft {
				middlewareRefs = append(middlewareRefs, mw.toRef())
			}
		case domain.AuthenticationTypeHard:
			for _, mw := range authConfig.Hard {
				middlewareRefs = append(middlewareRefs, mw.toRef())
			}
		}
	} else if website.Authentication != domain.AuthenticationTypeOff {
		log.Warnf("auth config not available for %s", website.FQDN)
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
	if b.useSablier(app) {
		middleware := b.sablierMiddleware(app)
		middlewares = append(middlewares, middleware)
		middlewareRefs = append(middlewareRefs, traefikv1alpha1.MiddlewareRef{Name: middleware.Name})
	}
	var rulePriority int
	{
		priorityOffset := b.config.Routing.Traefik.PriorityOffset
		rulePriority = len(rule) + priorityOffset
	}

	var tls *traefikv1alpha1.TLS
	if website.HTTPS {
		switch b.config.TLS.Type {
		case tlsTypeTraefik:
			targetDomain := b.config.TLS.Traefik.Wildcard.Domains.TLSTargetDomain(website)
			tls = &traefikv1alpha1.TLS{
				SecretName:   tlsSecretName(targetDomain),
				CertResolver: b.config.TLS.Traefik.CertResolver,
				Domains:      []types.Domain{{Main: targetDomain}},
			}
		case tlsTypeCertManager:
			targetDomain := b.config.TLS.CertManager.Wildcard.Domains.TLSTargetDomain(website)
			tls = &traefikv1alpha1.TLS{
				SecretName: tlsSecretName(targetDomain),
			}
		}
	}

	ingressRoute := &traefikv1alpha1.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressRoute",
			APIVersion: "traefik.io/v1alpha1",
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
				Priority:    rulePriority,
				Kind:        "Rule",
				Services:    serviceRefs,
				Middlewares: middlewareRefs,
			}},
			TLS: tls,
		},
	}

	return ingressRoute, middlewares
}

func (b *Backend) websiteCertificates(website *domain.Website) []*certmanagerv1.Certificate {
	var certs []*certmanagerv1.Certificate
	if website.HTTPS && b.config.TLS.Type == tlsTypeCertManager {
		targetDomain := b.config.TLS.CertManager.Wildcard.Domains.TLSTargetDomain(website)
		certs = append(certs, b.certificate(targetDomain))
	}
	return certs
}
