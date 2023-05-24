package k8simpl

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) runtimeSpec(app *domain.RuntimeDesiredState) (*appsv1.StatefulSet, *v1.Secret) {
	var secret *v1.Secret
	if len(app.Envs) > 0 {
		secret = &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName(app.App.ID),
				Namespace: b.config.Namespace,
				Labels:    b.appLabel(app.App.ID),
			},
			StringData: app.Envs,
		}
	}
	envs := lo.MapToSlice(app.Envs, func(key string, value string) v1.EnvVar {
		return v1.EnvVar{Name: key, ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
			LocalObjectReference: v1.LocalObjectReference{Name: deploymentName(app.App.ID)},
			Key:                  key,
		}}}
	})
	// make sure computed result is stable
	slices.SortFunc(envs, func(a, b v1.EnvVar) bool { return strings.Compare(a.Name, b.Name) < 0 })

	cont := v1.Container{
		Name:            podContainerName,
		Image:           app.ImageName + ":" + app.ImageTag,
		Env:             envs,
		Resources:       b.config.resourceRequirements(),
		ImagePullPolicy: v1.PullAlways,
	}
	if args := app.App.Config.BuildConfig.EntrypointArgs(); len(args) > 0 {
		cont.Command = args
	}
	if args := app.App.Config.BuildConfig.CommandArgs(); len(args) > 0 {
		cont.Args = args
	}
	slices.SortFunc(app.App.Websites, (*domain.Website).Compare)
	for _, website := range app.App.Websites {
		cont.Ports = append(cont.Ports, v1.ContainerPort{
			ContainerPort: int32(website.HTTPPort),
			Protocol:      v1.ProtocolTCP,
		})
	}
	slices.SortFunc(app.App.PortPublications, (*domain.PortPublication).Compare)
	for _, p := range app.App.PortPublications {
		cont.Ports = append(cont.Ports, v1.ContainerPort{
			ContainerPort: int32(p.ApplicationPort),
			Protocol:      protocolMapper.IntoMust(p.Protocol),
		})
	}
	cont.Ports = lo.UniqBy(cont.Ports, func(port v1.ContainerPort) string { return fmt.Sprintf("%d/%s", port.ContainerPort, port.Protocol) })
	ss := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(app.App.ID),
			Namespace: b.config.Namespace,
			Labels:    b.appLabel(app.App.ID),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: lo.ToPtr(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: appSelector(app.App.ID),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.appLabel(app.App.ID),
					Annotations: map[string]string{
						appRestartAnnotation: app.App.UpdatedAt.Format(time.RFC3339Nano),
					},
				},
				Spec: v1.PodSpec{
					AutomountServiceAccountToken: lo.ToPtr(false),
					EnableServiceLinks:           lo.ToPtr(false),
					Containers:                   []v1.Container{cont},
					NodeSelector:                 b.config.podSchedulingNodeSelector(),
					Tolerations:                  b.config.podSchedulingTolerations(),
				},
			},
		},
	}

	if b.config.ImagePullSecret != "" {
		ss.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{{Name: b.config.ImagePullSecret}}
	}

	return ss, secret
}

func (b *k8sBackend) runtimeResources(next *resources, apps []*domain.RuntimeDesiredState) {
	for _, app := range apps {
		ss, secret := b.runtimeSpec(app)
		next.statefulSets = append(next.statefulSets, ss)
		if secret != nil {
			next.secrets = append(next.secrets, secret)
		}
		for _, website := range app.App.Websites {
			next.services = append(next.services, b.runtimeService(app.App, website))
			ingressRoute, mw, certs := b.ingressRoute(app.App, website, b.runtimeServiceRef(app.App, website))
			next.middlewares = append(next.middlewares, mw...)
			next.ingressRoutes = append(next.ingressRoutes, ingressRoute)
			next.certificates = append(next.certificates, certs...)
		}
		for _, p := range app.App.PortPublications {
			next.services = append(next.services, b.runtimePortService(app.App, p))
		}
	}
}
