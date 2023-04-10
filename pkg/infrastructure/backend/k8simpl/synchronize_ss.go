package k8simpl

import (
	"context"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *k8sBackend) ssServiceRef() []traefikv1alpha1.Service {
	return []traefikv1alpha1.Service{{
		LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
			Name:      b.config.SS.Name,
			Kind:      b.config.SS.Kind,
			Namespace: b.config.SS.Namespace,
			Port:      intstr.FromInt(b.config.SS.Port),
			Scheme:    "http",
		},
	}}
}

func (b *k8sBackend) ssHeaderMiddleware(ss *domain.StaticSite) *traefikv1alpha1.Middleware {
	return &traefikv1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Middleware",
			APIVersion: "traefik.containo.us/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ssHeaderMiddlewareName(ss),
			Namespace: b.config.Namespace,
			Labels:    ssResourceLabels(ss.Application.ID),
		},
		Spec: traefikv1alpha1.MiddlewareSpec{
			Headers: &dynamic.Headers{
				CustomRequestHeaders: map[string]string{
					web.HeaderNameSSGenAppID: ss.Application.ID,
				},
			},
		},
	}
}

type ssResources struct {
	middlewares   []*traefikv1alpha1.Middleware
	ingressRoutes []*traefikv1alpha1.IngressRoute
	certificates  []*certmanagerv1.Certificate
}

func (b *k8sBackend) listCurrentSSResources(ctx context.Context) (*ssResources, error) {
	var resources ssResources
	listOpt := metav1.ListOptions{LabelSelector: ssLabelSelector()}

	mw, err := b.traefikClient.Middlewares(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get middlewares")
	}
	resources.middlewares = util.SliceOfPtr(mw.Items)

	ir, err := b.traefikClient.IngressRoutes(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ingress routes")
	}
	resources.ingressRoutes = util.SliceOfPtr(ir.Items)

	if b.config.TLS.Type == tlsTypeCertManager {
		certs, err := b.certManagerClient.CertmanagerV1().Certificates(b.config.Namespace).List(ctx, listOpt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get certificates")
		}
		resources.certificates = util.SliceOfPtr(certs.Items)
	}

	return &resources, nil
}

func (b *k8sBackend) SynchronizeSSIngress(ctx context.Context, sites []*domain.StaticSite, ads domain.AvailableDomainSlice) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	// List old resources
	old, err := b.listCurrentSSResources(ctx)
	if err != nil {
		return err
	}

	// Calculate next resources to apply
	var next ssResources
	for _, site := range sites {
		ingressRoute, mw, certs := b.ingressRoute(site.Application, site.Website, ssResourceLabels(site.Application.ID), b.ssServiceRef(), ads)

		ssHeaderMW := b.ssHeaderMiddleware(site)
		ingressRoute.Spec.Routes[0].Middlewares = append(ingressRoute.Spec.Routes[0].Middlewares, traefikv1alpha1.MiddlewareRef{Name: ssHeaderMW.Name})
		mw = append(mw, ssHeaderMW)

		next.middlewares = append(next.middlewares, mw...)
		next.ingressRoutes = append(next.ingressRoutes, ingressRoute)
		next.certificates = append(next.certificates, certs...)
	}
	next.certificates = lo.FindDuplicatesBy(next.certificates, func(cert *certmanagerv1.Certificate) string { return cert.Name })

	// Apply resources
	for _, mw := range next.middlewares {
		err = patch[*traefikv1alpha1.Middleware](ctx, mw.Name, mw, b.traefikClient.Middlewares(b.config.Namespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch middleware")
		}
	}
	for _, ir := range next.ingressRoutes {
		err = patch[*traefikv1alpha1.IngressRoute](ctx, ir.Name, ir, b.traefikClient.IngressRoutes(b.config.Namespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch ingress route")
		}
	}
	for _, cert := range next.certificates {
		err = patch[*certmanagerv1.Certificate](ctx, cert.Name, cert, b.certManagerClient.CertmanagerV1().Certificates(b.config.Namespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch certificate")
		}
	}

	// Prune old resources
	err = prune[*traefikv1alpha1.Middleware](ctx, diff(old.middlewares, next.middlewares), b.traefikClient.Middlewares(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune middlewares")
	}
	err = prune[*traefikv1alpha1.IngressRoute](ctx, diff(old.ingressRoutes, next.ingressRoutes), b.traefikClient.IngressRoutes(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune ingress route")
	}
	err = prune[*certmanagerv1.Certificate](ctx, diff(old.certificates, next.certificates), b.certManagerClient.CertmanagerV1().Certificates(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune certificates")
	}

	return nil
}
