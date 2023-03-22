package k8simpl

import (
	"context"
	"fmt"

	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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
		return fmt.Errorf("failed to get middlewares: %w", err)
	}
	existingIngressRoutes, err := b.traefikClient.IngressRoutes(appNamespace).List(ctx, listOpt)
	if err != nil {
		return fmt.Errorf("failed to get IngressRotues: %w", err)
	}

	// Calculate next resources to apply
	var middlewares []*v1alpha1.Middleware
	var ingressRoutes []*v1alpha1.IngressRoute
	for _, site := range sites {
		ingressRoute, mw := ingressRouteBase(site.Application, site.Website, ssResourceLabels(site.Application.ID))
		ingressRoute.Spec.Routes[0].Services = b.ssServiceRef()
		middlewares = append(middlewares, mw...)
		ingressRoutes = append(ingressRoutes, ingressRoute)
	}

	// Apply resources
	for _, mw := range middlewares {
		err = patch(ctx, mw.Name, mw, b.traefikClient.Middlewares(appNamespace))
		if err != nil {
			return fmt.Errorf("failed to patch middleware: %w", err)
		}
	}
	for _, ir := range ingressRoutes {
		err = patch(ctx, ir.Name, ir, b.traefikClient.IngressRoutes(appNamespace))
		if err != nil {
			return fmt.Errorf("failed to patch IngressRoute: %w", err)
		}
	}

	// Prune old resources
	err = prune(ctx, diff(util.SliceOfPtr(existingMiddlewares.Items), middlewares), b.traefikClient.Middlewares(appNamespace))
	if err != nil {
		return fmt.Errorf("failed to prune middlewares: %w", err)
	}
	err = prune(ctx, diff(util.SliceOfPtr(existingIngressRoutes.Items), ingressRoutes), b.traefikClient.IngressRoutes(appNamespace))
	if err != nil {
		return fmt.Errorf("failed to prune IngressRoutes: %w", err)
	}

	return nil
}
