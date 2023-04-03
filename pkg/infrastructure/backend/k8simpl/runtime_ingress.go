package k8simpl

import (
	"context"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	"github.com/traefik/traefik/v2/pkg/types"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util"
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

func ingressRouteBase(app *domain.Application, website *domain.Website, labels map[string]string) (*traefikv1alpha1.IngressRoute, []*traefikv1alpha1.Middleware) {
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
				Middlewares: middlewareRefs,
			}},
			TLS: tls,
		},
	}

	return ingressRoute, middlewares
}

func (b *k8sBackend) listRuntimeIngressResources(ctx context.Context, appID string) (services *apiv1.ServiceList, middlewares *traefikv1alpha1.MiddlewareList, ingressRoutes *traefikv1alpha1.IngressRouteList, err error) {
	listOpt := metav1.ListOptions{LabelSelector: labelSelector(appID)}
	services, err = b.client.CoreV1().Services(appNamespace).List(ctx, listOpt)
	if err != nil {
		err = errors.Wrap(err, "failed to get services")
		return
	}
	middlewares, err = b.traefikClient.Middlewares(appNamespace).List(ctx, listOpt)
	if err != nil {
		err = errors.Wrap(err, "failed to get middlewares")
		return
	}
	ingressRoutes, err = b.traefikClient.IngressRoutes(appNamespace).List(ctx, listOpt)
	if err != nil {
		err = errors.Wrap(err, "failed to get IngressRotues")
		return
	}
	return
}

func (b *k8sBackend) synchronizeRuntimeIngresses(ctx context.Context, app *domain.Application) error {
	// Collect current resources
	existingServices, existingMiddlewares, existingIngressRoutes, err := b.listRuntimeIngressResources(ctx, app.ID)
	if err != nil {
		return err
	}

	// Calculate next resources to apply
	var services []*apiv1.Service
	var middlewares []*traefikv1alpha1.Middleware
	var ingressRoutes []*traefikv1alpha1.IngressRoute
	for _, website := range app.Websites {
		services = append(services, runtimeService(app, website))
		ingressRoute, mw := ingressRouteBase(app, website, resourceLabels(app.ID))
		ingressRoute.Spec.Routes[0].Services = runtimeServiceRef(app, website)
		middlewares = append(middlewares, mw...)
		ingressRoutes = append(ingressRoutes, ingressRoute)
	}

	// Apply resources
	for _, svc := range services {
		err = patch[*apiv1.Service](ctx, svc.Name, svc, b.client.CoreV1().Services(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch service")
		}
	}
	for _, mw := range middlewares {
		err = patch[*traefikv1alpha1.Middleware](ctx, mw.Name, mw, b.traefikClient.Middlewares(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch middleware")
		}
	}
	for _, ir := range ingressRoutes {
		err = patch[*traefikv1alpha1.IngressRoute](ctx, ir.Name, ir, b.traefikClient.IngressRoutes(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch IngressRoute")
		}
	}

	// Prune old resources
	err = prune(ctx, diff(util.SliceOfPtr(existingServices.Items), services), b.client.CoreV1().Services(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune services")
	}
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

func (b *k8sBackend) destroyRuntimeIngresses(ctx context.Context, app *domain.Application) error {
	services, middlewares, ingressRoutes, err := b.listRuntimeIngressResources(ctx, app.ID)
	if err != nil {
		return err
	}

	err = prune(ctx, names(util.SliceOfPtr(services.Items)), b.client.CoreV1().Services(appNamespace))
	if err != nil {
		return err
	}
	err = prune(ctx, names(util.SliceOfPtr(middlewares.Items)), b.traefikClient.Middlewares(appNamespace))
	if err != nil {
		return err
	}
	err = prune(ctx, names(util.SliceOfPtr(ingressRoutes.Items)), b.traefikClient.IngressRoutes(appNamespace))
	if err != nil {
		return err
	}

	return nil
}
