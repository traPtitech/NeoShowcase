package k8simpl

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *k8sBackend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
	list, err := b.client.CoreV1().Pods(b.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: toSelectorString(appSelector(appID)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pods")
	}

	if len(list.Items) == 0 {
		return &domain.Container{
			ApplicationID: appID,
			State:         domain.ContainerStateMissing,
		}, nil
	}
	return &domain.Container{
		ApplicationID: appID,
		State:         getContainerState(list.Items[0].Status),
	}, nil
}

func (b *k8sBackend) ListContainers(ctx context.Context) ([]*domain.Container, error) {
	list, err := b.client.CoreV1().Pods(b.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: toSelectorString(allSelector()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pods")
	}

	result := lo.Map(list.Items, func(pod v1.Pod, i int) *domain.Container {
		return &domain.Container{
			ApplicationID: pod.Labels[appIDLabel],
			State:         getContainerState(pod.Status),
		}
	})
	return result, nil
}

func getContainerState(status v1.PodStatus) domain.ContainerState {
	switch status.Phase {
	case v1.PodPending:
		return domain.ContainerStateStarting
	case v1.PodRunning:
		return domain.ContainerStateRunning
	case v1.PodFailed:
		return domain.ContainerStateErrored
	case v1.PodSucceeded:
		return domain.ContainerStateExited
	default:
		return domain.ContainerStateUnknown
	}
}
