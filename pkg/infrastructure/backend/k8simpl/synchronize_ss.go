package k8simpl

import (
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
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
			Labels:    b.appLabel(ss.Application.ID),
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

func (b *k8sBackend) ssResources(next *resources, sites []*domain.StaticSite, ads domain.AvailableDomainSlice) {
	for _, site := range sites {
		ingressRoute, mw, certs := b.ingressRoute(site.Application, site.Website, b.ssServiceRef(), ads)

		ssHeaderMW := b.ssHeaderMiddleware(site)
		ingressRoute.Spec.Routes[0].Middlewares = append(ingressRoute.Spec.Routes[0].Middlewares, traefikv1alpha1.MiddlewareRef{Name: ssHeaderMW.Name})
		mw = append(mw, ssHeaderMW)

		next.middlewares = append(next.middlewares, mw...)
		next.ingressRoutes = append(next.ingressRoutes, ingressRoute)
		next.certificates = append(next.certificates, certs...)
	}
}
