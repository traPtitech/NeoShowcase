package k8simpl

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (m *Manager) Create(ctx context.Context, args container.CreateArgs) (*container.CreateResult, error) {
	if args.ImageTag == "" {
		args.ImageTag = "latest"
	}

	labels := util.MergeLabels(args.Labels, map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: args.ApplicationID,
		appContainerEnvironmentIDLabel: args.EnvironmentID,
	})

	var envs []apiv1.EnvVar

	for name, value := range args.Envs {
		envs = append(envs, apiv1.EnvVar{Name: name, Value: value})
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
				Name:      deploymentName(args.ApplicationID, args.EnvironmentID),
				Namespace: appNamespace,
				Labels:    labels,
			},
			Spec: apiv1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					appContainerLabel:              "true",
					appContainerApplicationIDLabel: args.ApplicationID,
					appContainerEnvironmentIDLabel: args.EnvironmentID,
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
				Name:      deploymentName(args.ApplicationID, args.EnvironmentID),
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
											ServiceName: deploymentName(args.ApplicationID, args.EnvironmentID),
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
			if args.Recreate && errors.IsAlreadyExists(err) {
				if err = m.clientset.CoreV1().Services(appNamespace).Delete(ctx, svc.Name, metav1.DeleteOptions{}); err != nil {
					return nil, fmt.Errorf("failed to delete service: %w", err)
				}
				if _, err := m.clientset.CoreV1().Services(appNamespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
					return nil, fmt.Errorf("failed to create service: %w", err)
				}
			} else {
				return nil, fmt.Errorf("failed to create service: %w", err)
			}
		}
		if _, err := m.clientset.NetworkingV1beta1().Ingresses(appNamespace).Create(ctx, ingress, metav1.CreateOptions{}); err != nil {
			if args.Recreate && errors.IsAlreadyExists(err) {
				if _, err = m.clientset.NetworkingV1beta1().Ingresses(appNamespace).Update(ctx, ingress, metav1.UpdateOptions{}); err != nil {
					return nil, fmt.Errorf("failed to update ingress: %w", err)
				}
			} else {
				return nil, fmt.Errorf("failed to create ingress: %w", err)
			}
		}
	}

	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(args.ApplicationID, args.EnvironmentID),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{cont},
		},
	}

	_, err := m.clientset.CoreV1().Pods(appNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if args.Recreate && errors.IsAlreadyExists(err) {
			// TODO いい感じにしたい
			if err = m.clientset.CoreV1().Pods(appNamespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
				return nil, fmt.Errorf("failed to delete pod: %w", err)
			}
			for {
				_, err := m.clientset.CoreV1().Pods(appNamespace).Get(ctx, pod.Name, metav1.GetOptions{})
				if err != nil {
					if errors.IsNotFound(err) {
						break
					}
					return nil, fmt.Errorf("failed to delete pod: %w", err)
				}
				time.Sleep(500 * time.Millisecond)
			}
			if _, err := m.clientset.CoreV1().Pods(appNamespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
				return nil, fmt.Errorf("failed to create pod: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to create pod: %w", err)
		}
	}

	return &container.CreateResult{}, nil
}
