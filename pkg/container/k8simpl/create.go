package k8simpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (m *Manager) Create(ctx context.Context, args container.CreateArgs) (*container.CreateResult, error) {
	labels := util.MergeLabels(args.Labels, map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: args.ApplicationID,
	})

	var envs []apiv1.EnvVar

	for _, env := range args.Envs {
		envs = append(envs, apiv1.EnvVar{Name: strings.SplitN(env, "=", 2)[0], Value: strings.SplitN(env, "=", 2)[1]})
	}

	cont := apiv1.Container{
		Name:  "app",
		Image: args.ImageName + ":" + args.ImageTag,
		Env:   envs,
	}
	if args.HTTPProxy != nil {
		cont.Ports = []apiv1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: int32(args.HTTPProxy.Port),
				Protocol:      "TCP",
			},
		}
		svc := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName(args.ApplicationID),
				Namespace: appNamespace,
				Labels:    labels,
			},
			Spec: apiv1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					appContainerLabel:              "true",
					appContainerApplicationIDLabel: args.ApplicationID,
				},
				Ports: []apiv1.ServicePort{
					{
						Protocol:   "TCP",
						Port:       80,
						TargetPort: intstr.FromString("http"),
					},
				},
			},
		}
		ingress := &networkingv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName(args.ApplicationID),
				Namespace: appNamespace,
				Labels:    labels,
				Annotations: map[string]string{
					"traefik.ingress.kubernetes.io/router.entrypoints": "web, websecure", // TODO HTTPS可能かどうかを判断
				},
			},
			Spec: networkingv1beta1.IngressSpec{
				Rules: []networkingv1beta1.IngressRule{
					{
						Host: args.HTTPProxy.Domain,
						IngressRuleValue: networkingv1beta1.IngressRuleValue{
							HTTP: &networkingv1beta1.HTTPIngressRuleValue{
								Paths: []networkingv1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: networkingv1beta1.IngressBackend{
											ServiceName: deploymentName(args.ApplicationID),
											ServicePort: intstr.FromInt(80),
										},
									},
								},
							},
						},
					},
				},
			},
		}

		if _, err := m.clientset.CoreV1().Services(appNamespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create service: %w", err)
		}
		if _, err := m.clientset.NetworkingV1beta1().Ingresses(appNamespace).Create(ctx, ingress, metav1.CreateOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create ingress: %w", err)
		}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(args.ApplicationID),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					appContainerLabel:              "true",
					appContainerApplicationIDLabel: args.ApplicationID,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: appNamespace,
					Labels:    labels,
				},
				Spec: apiv1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers:    []apiv1.Container{cont},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
		},
	}

	_, err := m.clientset.AppsV1().Deployments(appNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	return &container.CreateResult{}, nil
}
