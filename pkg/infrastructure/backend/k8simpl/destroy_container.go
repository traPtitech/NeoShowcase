package k8simpl

import (
	"context"

	"github.com/friendsofgo/errors"
	v1 "k8s.io/api/apps/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type (
	m map[string]any
)

func (b *k8sBackend) DestroyContainer(ctx context.Context, app *domain.Application) error {
	err := b.destroyRuntimeIngresses(ctx, app)
	if err != nil {
		return errors.Wrap(err, "failed to destroy runtime ingress resources")
	}

	// statefulset の spec.selector がなぜか omitempty ではないため
	statefulSetName := deploymentName(app.ID)
	statefulSet := m{
		"kind":       "StatefulSet",
		"apiVersion": "apps/v1",
		"metadata": m{
			"name":      statefulSetName,
			"namespace": appNamespace,
		},
		"spec": m{
			"replicas": 0,
		},
	}

	return strategicPatch[*v1.StatefulSet](ctx, statefulSetName, statefulSet, b.client.AppsV1().StatefulSets(appNamespace))
}
