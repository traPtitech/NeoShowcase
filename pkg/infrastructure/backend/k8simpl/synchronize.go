package k8simpl

import (
	"context"
	"time"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikcontainous/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

type runtimeResources struct {
	statefulSets  []*appsv1.StatefulSet
	services      []*v1.Service
	middlewares   []*traefikv1alpha1.Middleware
	ingressRoutes []*traefikv1alpha1.IngressRoute
	certificates  []*certmanagerv1.Certificate
}

func (b *k8sBackend) listCurrentResources(ctx context.Context) (*runtimeResources, error) {
	var resources runtimeResources
	listOpt := metav1.ListOptions{LabelSelector: allSelector()}

	ss, err := b.client.AppsV1().StatefulSets(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stateful sets")
	}
	resources.statefulSets = util.SliceOfPtr(ss.Items)

	svc, err := b.client.CoreV1().Services(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get services")
	}
	resources.services = util.SliceOfPtr(svc.Items)

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

func (b *k8sBackend) statefulSet(app *domain.AppDesiredState) *appsv1.StatefulSet {
	envs := lo.MapToSlice(app.Envs, func(key string, value string) v1.EnvVar {
		return v1.EnvVar{Name: key, Value: value}
	})

	cont := v1.Container{
		Name:  "app",
		Image: app.ImageName + ":" + app.ImageTag,
		Env:   envs,
	}
	if app.App.Config.Entrypoint != "" {
		cont.Command = app.App.Config.EntrypointArgs()
	}
	if app.App.Config.Command != "" {
		cont.Args = app.App.Config.CommandArgs()
	}
	ss := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(app.App.ID),
			Namespace: b.config.Namespace,
			Labels:    resourceLabels(app.App.ID),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: lo.ToPtr(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: resourceLabels(app.App.ID),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: resourceLabels(app.App.ID),
					Annotations: map[string]string{
						appRestartAnnotation: app.App.UpdatedAt.Format(time.RFC3339),
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{cont},
				},
			},
		},
	}

	if b.config.ImagePullSecret != "" {
		ss.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{{Name: b.config.ImagePullSecret}}
	}

	return ss
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
		next.statefulSets = append(next.statefulSets, b.statefulSet(app))
		for _, website := range app.App.Websites {
			next.services = append(next.services, b.runtimeService(app.App, website))
			ingressRoute, mw, certs := b.ingressRoute(app.App, website, resourceLabels(app.App.ID), b.runtimeServiceRef(app.App, website))
			next.middlewares = append(next.middlewares, mw...)
			next.ingressRoutes = append(next.ingressRoutes, ingressRoute)
			next.certificates = append(next.certificates, certs...)
		}
	}
	next.certificates = lo.FindDuplicatesBy(next.certificates, func(cert *certmanagerv1.Certificate) string { return cert.Name })

	// Apply resources
	for _, ss := range next.statefulSets {
		err = patch[*appsv1.StatefulSet](ctx, ss.Name, ss, b.client.AppsV1().StatefulSets(b.config.Namespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch stateful set")
		}
	}
	for _, svc := range next.services {
		err = patch[*v1.Service](ctx, svc.Name, svc, b.client.CoreV1().Services(b.config.Namespace))
		if err != nil {
			return errors.Wrap(err, "failed to patch service")
		}
	}
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
	// NOTE: stateful set does not provide any guarantees to (order of) termination of pods when deleted
	// NOTE: stateful set does not delete volumes
	// https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#limitations
	err = prune[*appsv1.StatefulSet](ctx, diff(old.statefulSets, next.statefulSets), b.client.AppsV1().StatefulSets(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune stateful sets")
	}
	err = prune[*v1.Service](ctx, diff(old.services, next.services), b.client.CoreV1().Services(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune services")
	}
	err = prune[*traefikv1alpha1.Middleware](ctx, diff(old.middlewares, next.middlewares), b.traefikClient.Middlewares(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune middlewares")
	}
	err = prune[*traefikv1alpha1.IngressRoute](ctx, diff(old.ingressRoutes, next.ingressRoutes), b.traefikClient.IngressRoutes(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune ingress routes")
	}
	err = prune[*certmanagerv1.Certificate](ctx, diff(old.certificates, next.certificates), b.certManagerClient.CertmanagerV1().Certificates(b.config.Namespace))
	if err != nil {
		return errors.Wrap(err, "failed to prune certificates")
	}

	return nil
}
