package k8simpl

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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

	result := ds.Map(list.Items, func(pod v1.Pod) *domain.Container {
		return &domain.Container{
			ApplicationID: pod.Labels[appIDLabel],
			State:         getContainerState(pod.Status),
		}
	})
	return result, nil
}

func getContainerState(status v1.PodStatus) domain.ContainerState {
	cs, ok := lo.Find(status.ContainerStatuses, func(cs v1.ContainerStatus) bool { return cs.Name == podContainerName })
	if !ok {
		return domain.ContainerStateMissing
	}
	if cs.State.Waiting != nil {
		return domain.ContainerStateStarting
	}
	if cs.State.Running != nil {
		return domain.ContainerStateRunning
	}
	if cs.State.Terminated != nil {
		if cs.State.Terminated.ExitCode == 0 {
			return domain.ContainerStateExited
		} else {
			return domain.ContainerStateErrored
		}
	}
	return domain.ContainerStateUnknown
}
