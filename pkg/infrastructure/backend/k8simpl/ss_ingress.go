package k8simpl

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *k8sBackend) ssServiceRef() []v1alpha1.Service {
	return []v1alpha1.Service{{
		LoadBalancerSpec: v1alpha1.LoadBalancerSpec{
			Name:      b.ss.Service.Name,
			Kind:      b.ss.Service.Kind,
			Namespace: b.ss.Service.Namespace,
			Port:      intstr.FromInt(b.ss.Service.Port),
			Scheme:    "http",
		},
	}}
}

func ssHeaderMiddleware(ss *domain.StaticSite) *v1alpha1.Middleware {
	return &v1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Middleware",
			APIVersion: "traefik.containo.us/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ssHeaderMiddlewareName(ss),
			Namespace: appNamespace,
			Labels:    ssResourceLabels(ss.Application.ID),
		},
		Spec: v1alpha1.MiddlewareSpec{
			Headers: &dynamic.Headers{
				CustomRequestHeaders: map[string]string{
					web.HeaderNameSSGenAppID: ss.Application.ID,
				},
			},
		},
	}
}

func (b *k8sBackend) ReloadSSIngress(ctx context.Context) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	sites, err := domain.GetActiveStaticSites(context.Background(), b.appRepo, b.buildRepo)
	if err != nil {
		return err
	}

	// Collect current resources
	listOpt := metav1.ListOptions{LabelSelector: ssLabelSelector()}
	existingMiddlewares, err := b.traefikClient.Middlewares(appNamespace).List(ctx, listOpt)
	if err != nil {
		return errors.Wrap(err, "failed to get middlewares")
	}
	existingIngressRoutes, err := b.traefikClient.IngressRoutes(appNamespace).List(ctx, listOpt)
	if err != nil {
		return errors.Wrap(err, "failed to get IngressRotues")
	}

	// Calculate next resources to apply
	var middlewares []*v1alpha1.Middleware
	var ingressRoutes []*v1alpha1.IngressRoute
	for _, site := range sites {
		ingressRoute, mw := ingressRouteBase(site.Application, site.Website, ssResourceLabels(site.Application.ID))
		ingressRoute.Spec.Routes[0].Services = b.ssServiceRef()

		ssHeaderMW := ssHeaderMiddleware(site)
		ingressRoute.Spec.Routes[0].Middlewares = append(ingressRoute.Spec.Routes[0].Middlewares, v1alpha1.MiddlewareRef{Name: ssHeaderMW.Name})
		mw = append(mw, ssHeaderMW)

		middlewares = append(middlewares, mw...)
		ingressRoutes = append(ingressRoutes, ingressRoute)
	}

	// Apply resources
	for _, mw := range middlewares {
		err = patch[*v1alpha1.Middleware](ctx, mw.Name, mw, b.traefikClient.Middlewares(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch middleware")
		}
	}
	for _, ir := range ingressRoutes {
		err = patch[*v1alpha1.IngressRoute](ctx, ir.Name, ir, b.traefikClient.IngressRoutes(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch IngressRoute")
		}
	}

	// Prune old resources
	err = prune(ctx, diff(util.SliceOfPtr(existingMiddlewares.Items), middlewares), b.traefikClient.Middlewares(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune middlewares")
	}
	err = prune(ctx, diff(util.SliceOfPtr(existingIngressRoutes.Items), ingressRoutes), b.traefikClient.IngressRoutes(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune IngressRoutes")
	}

	return nil
}