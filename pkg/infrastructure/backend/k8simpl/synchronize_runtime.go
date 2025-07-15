package k8simpl

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

func comparePort(port v1.ContainerPort) string {
	return fmt.Sprintf("%d/%s", port.ContainerPort, port.Protocol)
}

func (b *Backend) runtimeSpec(app *domain.RuntimeDesiredState) (*appsv1.StatefulSet, *v1.Service, *v1.Secret) {
	var secret *v1.Secret
	var envs []v1.EnvVar
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
		// NOTE: marshaling map[string]string is stable (json.Marshal sorts by key)
		b, _ := json.Marshal(secret)
		secret.Name = deploymentNameWithDiscriminator(app.App.ID, b)

		envs = lo.MapToSlice(app.Envs, func(key string, value string) v1.EnvVar {
			return v1.EnvVar{Name: key, ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{Name: secret.Name},
				Key:                  key,
			}}}
		})
		// make sure computed result is stable
		slices.SortFunc(envs, ds.LessFunc(func(a v1.EnvVar) string { return a.Name }))
	}

	cont := v1.Container{
		Name:            podContainerName,
		Image:           app.ImageName + ":" + app.ImageTag,
		Env:             envs,
		Resources:       b.config.resourceRequirements(),
		ImagePullPolicy: v1.PullAlways,
		Stdin:           true,
		TTY:             true,
	}
	if args, _ := domain.ParseArgs(app.App.Config.BuildConfig.GetRuntimeConfig().Entrypoint); len(args) > 0 {
		cont.Command = args
	}
	if args, _ := domain.ParseArgs(app.App.Config.BuildConfig.GetRuntimeConfig().Command); len(args) > 0 {
		cont.Args = args
	}

	for _, website := range app.App.Websites {
		cont.Ports = append(cont.Ports, v1.ContainerPort{
			ContainerPort: int32(website.HTTPPort),
			Protocol:      v1.ProtocolTCP,
		})
	}
	for _, p := range app.App.PortPublications {
		cont.Ports = append(cont.Ports, v1.ContainerPort{
			ContainerPort: int32(p.ApplicationPort),
			Protocol:      protocolMapper.IntoMust(p.Protocol),
		})
	}
	cont.Ports = lo.UniqBy(cont.Ports, comparePort)
	slices.SortFunc(cont.Ports, ds.LessFunc(comparePort))

	var replicas = int32(1)
	var ssLabels = b.appLabel(app.App.ID)

	if b.useSablier(app.App) {
		ssLabels["sablier.enable"] = "true"
		ssLabels["sablier.group"] = sablierGroupName(app.App.ID)
		replicas = int32(0) // scaled by sablier
	}

	ss := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(app.App.ID),
			Namespace: b.config.Namespace,
			Labels:    ssLabels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: appSelector(app.App.ID),
			},
			// to not wait for Pods to become Running and Ready or completely terminated prior to launching or terminating another Pod
			// https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#parallel-pod-management
			PodManagementPolicy: appsv1.ParallelPodManagement,
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
					NodeSelector:                 b.config.podSchedulingNodeSelector(app.App.ID),
					Tolerations:                  b.config.podSchedulingTolerations(),
					TopologySpreadConstraints:    b.config.podSpreadConstraints(),
				},
			},
			RevisionHistoryLimit: lo.ToPtr(int32(0)),
		},
	}

	if b.config.ImagePullSecret != "" {
		ss.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{{Name: b.config.ImagePullSecret}}
	}

	var svc *v1.Service
	if len(cont.Ports) > 0 {
		svc = &v1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName(app.App.ID),
				Namespace: b.config.Namespace,
				Labels:    b.appLabel(app.App.ID),
			},
			Spec: v1.ServiceSpec{
				Type:           "ClusterIP",
				IPFamilies:     b.config.serviceIPFamilies(),
				IPFamilyPolicy: b.config.serviceIPFamilyPolicy(),
				Selector:       appSelector(app.App.ID),
				Ports: ds.Map(cont.Ports, func(port v1.ContainerPort) v1.ServicePort {
					return v1.ServicePort{
						Name:       fmt.Sprintf("%v-%v", strings.ToLower(string(port.Protocol)), port.ContainerPort),
						Protocol:   port.Protocol,
						Port:       port.ContainerPort,
						TargetPort: intstr.FromInt(int(port.ContainerPort)),
					}
				}),
			},
		}
	}

	return ss, svc, secret
}

func (b *Backend) runtimeServiceRef(app *domain.Application, website *domain.Website) []traefikv1alpha1.Service {
	return []traefikv1alpha1.Service{{
		LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
			Name:      deploymentName(app.ID),
			Kind:      "Service",
			Namespace: b.config.Namespace,
			Port:      intstr.FromInt(website.HTTPPort),
			Scheme:    lo.Ternary(website.H2C, "h2c", "http"),
		},
	}}
}

var protocolMapper = mapper.MustNewValueMapper(map[domain.PortPublicationProtocol]v1.Protocol{
	domain.PortPublicationProtocolTCP: v1.ProtocolTCP,
	domain.PortPublicationProtocolUDP: v1.ProtocolUDP,
})

func (b *Backend) runtimePortService(app *domain.Application, port *domain.PortPublication) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      portServiceName(port),
			Namespace: b.config.Namespace,
			Labels:    b.appLabel(app.ID),
		},
		Spec: v1.ServiceSpec{
			Type:           "LoadBalancer",
			IPFamilies:     b.config.serviceIPFamilies(),
			IPFamilyPolicy: b.config.serviceIPFamilyPolicy(),
			Selector:       appSelector(app.ID),
			Ports: []v1.ServicePort{{
				Protocol:   protocolMapper.IntoMust(port.Protocol),
				Port:       int32(port.InternetPort),
				TargetPort: intstr.FromInt(port.ApplicationPort),
			}},
		},
	}
}

func (b *Backend) runtimeResources(next *resources, apps []*domain.RuntimeDesiredState) {
	for _, app := range apps {
		// Filter to sharded apps
		if !b.cluster.IsAssigned(app.App.ID) {
			continue
		}

		ss, svc, secret := b.runtimeSpec(app)
		next.statefulSets = append(next.statefulSets, ss)
		if svc != nil {
			next.services = append(next.services, svc)
		}
		if secret != nil {
			next.secrets = append(next.secrets, secret)
		}
		for _, website := range app.App.Websites {
			ingressRoute, mw := b.ingressRoute(app.App, website, b.runtimeServiceRef(app.App, website))
			next.middlewares = append(next.middlewares, mw...)
			next.ingressRoutes = append(next.ingressRoutes, ingressRoute)
		}
		for _, p := range app.App.PortPublications {
			next.services = append(next.services, b.runtimePortService(app.App, p))
		}
	}
}
