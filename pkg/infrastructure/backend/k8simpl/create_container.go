package k8simpl

import (
	"context"
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *k8sBackend) CreateContainer(ctx context.Context, args domain.ContainerCreateArgs) error {
	if args.ImageTag == "" {
		args.ImageTag = "latest"
	}

	labels := util.MergeLabels(args.Labels, map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: args.ApplicationID,
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

		if _, err := b.clientset.CoreV1().Services(appNamespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
			if args.Recreate && errors.IsAlreadyExists(err) {
				if err = b.clientset.CoreV1().Services(appNamespace).Delete(ctx, svc.Name, metav1.DeleteOptions{}); err != nil {
					return fmt.Errorf("failed to delete service: %w", err)
				}
				if _, err := b.clientset.CoreV1().Services(appNamespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
					return fmt.Errorf("failed to create service: %w", err)
				}
			} else {
				return fmt.Errorf("failed to create service: %w", err)
			}
		}

		// TODO: use traefik to expose service?
		// https://doc.traefik.io/traefik/reference/dynamic-configuration/kubernetes-crd/#resources
	}

	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(args.ApplicationID),
			Namespace: appNamespace,
			Labels:    labels,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{cont},
		},
	}

	_, err := b.clientset.CoreV1().Pods(appNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if args.Recreate && errors.IsAlreadyExists(err) {
			// TODO いい感じにしたい
			if err = b.clientset.CoreV1().Pods(appNamespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("failed to delete pod: %w", err)
			}
			for {
				_, err := b.clientset.CoreV1().Pods(appNamespace).Get(ctx, pod.Name, metav1.GetOptions{})
				if err != nil {
					if errors.IsNotFound(err) {
						break
					}
					return fmt.Errorf("failed to delete pod: %w", err)
				}
				time.Sleep(500 * time.Millisecond)
			}
			if _, err := b.clientset.CoreV1().Pods(appNamespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
				return fmt.Errorf("failed to create pod: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create pod: %w", err)
		}
	}

	return nil

}
