package k8simpl

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	err := b.destroyRuntimeIngresses(ctx, app)
	if err != nil {
		return fmt.Errorf("failed to destroy runtime ingress resources: %w", err)
	}

	statefulSet := &v1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(app.ID),
			Namespace: appNamespace,
		},
		Spec: v1.StatefulSetSpec{
			Replicas: pointer.Int32(0),
		},
	}

	return strategicPatch[*v1.StatefulSet](ctx, statefulSet.Name, statefulSet, b.client.AppsV1().StatefulSets(appNamespace))
}
