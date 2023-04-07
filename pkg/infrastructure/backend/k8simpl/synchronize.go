package k8simpl

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

type runtimeResources struct {
	statefulSets  []*v1.StatefulSet
	services      []*apiv1.Service
	middlewares   []*traefikv1alpha1.Middleware
	ingressRoutes []*traefikv1alpha1.IngressRoute
}

func (b *k8sBackend) listCurrentResources(ctx context.Context) (*runtimeResources, error) {
	var resources runtimeResources
	listOpt := metav1.ListOptions{LabelSelector: allSelector()}

	ss, err := b.client.AppsV1().StatefulSets(appNamespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stateful sets")
	}
	resources.statefulSets = util.SliceOfPtr(ss.Items)

	svc, err := b.client.CoreV1().Services(appNamespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get services")
	}
	resources.services = util.SliceOfPtr(svc.Items)

	mw, err := b.traefikClient.Middlewares(appNamespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get middlewares")
	}
	resources.middlewares = util.SliceOfPtr(mw.Items)

	ir, err := b.traefikClient.IngressRoutes(appNamespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ingress routes")
	}
	resources.ingressRoutes = util.SliceOfPtr(ir.Items)

	return &resources, nil
}

func statefulSet(app *domain.AppDesiredState) *v1.StatefulSet {
	envs := lo.MapToSlice(app.Envs, func(key string, value string) apiv1.EnvVar {
		return apiv1.EnvVar{Name: key, Value: value}
	})

	return &v1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(app.App.ID),
			Namespace: appNamespace,
			Labels:    resourceLabels(app.App.ID),
		},
		Spec: v1.StatefulSetSpec{
			Replicas: lo.ToPtr(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: resourceLabels(app.App.ID),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: resourceLabels(app.App.ID),
					Annotations: map[string]string{
						appRestartAnnotation: app.App.UpdatedAt.Format(time.RFC3339),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{{
						Name:  "app",
						Image: app.ImageName + ":" + app.ImageTag,
						Env:   envs,
					}},
				},
			},
		},
	}
}

func (b *k8sBackend) deleteStatefulSet(ctx context.Context, ss *v1.StatefulSet) error {
	// HACK?: statefulset の spec.selector がなぜか omitempty ではないため、生 json を指定する
	patch := m{
		"kind":       "StatefulSet",
		"apiVersion": "apps/v1",
		"metadata": m{
			"name":      ss.Name,
			"namespace": appNamespace,
		},
		"spec": m{
			"replicas": 0,
		},
	}
	return strategicPatch[*v1.StatefulSet](ctx, ss.Name, patch, b.client.AppsV1().StatefulSets(appNamespace))
}

func (b *k8sBackend) SynchronizeRuntime(ctx context.Context, apps []*domain.AppDesiredState) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	// List old resources
	old, err := b.listCurrentResources(ctx)
	if err != nil {
		return err
	}

	// Calculate next resources to apply
	var next runtimeResources
	for _, app := range apps {
		next.statefulSets = append(next.statefulSets, statefulSet(app))
		for _, website := range app.App.Websites {
			next.services = append(next.services, runtimeService(app.App, website))
			ingressRoute, mw := ingressRoute(app.App, website, resourceLabels(app.App.ID), runtimeServiceRef(app.App, website))
			next.middlewares = append(next.middlewares, mw...)
			next.ingressRoutes = append(next.ingressRoutes, ingressRoute)
		}
	}

	// Apply resources
	for _, ss := range next.statefulSets {
		err = patch[*v1.StatefulSet](ctx, ss.Name, ss, b.client.AppsV1().StatefulSets(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch stateful set")
		}
	}
	for _, svc := range next.services {
		err = patch[*apiv1.Service](ctx, svc.Name, svc, b.client.CoreV1().Services(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch service")
		}
	}
	for _, mw := range next.middlewares {
		err = patch[*traefikv1alpha1.Middleware](ctx, mw.Name, mw, b.traefikClient.Middlewares(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch middleware")
		}
	}
	for _, ir := range next.ingressRoutes {
		err = patch[*traefikv1alpha1.IngressRoute](ctx, ir.Name, ir, b.traefikClient.IngressRoutes(appNamespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch ingress route")
		}
	}

	// Prune old resources
	// NOTE: stateful set does not provide any guarantees to (order of) termination of pods when deleted
	// NOTE: stateful set does not delete volumes
	// https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#limitations
	err = prune[*v1.StatefulSet](ctx, diff(old.statefulSets, next.statefulSets), b.client.AppsV1().StatefulSets(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune stateful sets")
	}
	err = prune[*apiv1.Service](ctx, diff(old.services, next.services), b.client.CoreV1().Services(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune services")
	}
	err = prune[*traefikv1alpha1.Middleware](ctx, diff(old.middlewares, next.middlewares), b.traefikClient.Middlewares(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune middlewares")
	}
	err = prune[*traefikv1alpha1.IngressRoute](ctx, diff(old.ingressRoutes, next.ingressRoutes), b.traefikClient.IngressRoutes(appNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune ingress routes")
	}

	return nil
}
