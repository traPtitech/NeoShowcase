package k8simpl

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *k8sBackend) CreateContainer(ctx context.Context, app *domain.Application, args domain.ContainerCreateArgs) error {
	if args.ImageTag == "" {
		args.ImageTag = "latest"
	}

	envs := lo.MapToSlice(args.Envs, func(key string, value string) apiv1.EnvVar {
		return apiv1.EnvVar{Name: key, Value: value}
	})

	err := b.synchronizeRuntimeIngresses(ctx, app)
	if err != nil {
		return fmt.Errorf("failed to synchronize ingresses: %w", err)
	}

	statefulSet := &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(app.ID),
			Namespace: appNamespace,
			Labels:    resourceLabels(app.ID),
		},
		Spec: v1.StatefulSetSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: resourceLabels(app.ID),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName(app.ID),
					Namespace: appNamespace,
					Labels: util.MergeMap(resourceLabels(app.ID), map[string]string{
						appRestartAnnotation: time.Now().Format(time.RFC3339),
					}),
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{{
						Name:  "app",
						Image: args.ImageName + ":" + args.ImageTag,
						Env:   envs,
					}},
				},
			},
		},
	}

	return patch(ctx, statefulSet.Name, statefulSet, b.client.AppsV1().StatefulSets(appNamespace))
}
