package k8simpl

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/fmtutil"
)

func (b *Backend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
	sts, err := b.client.AppsV1().StatefulSets(b.config.Namespace).Get(ctx, deploymentName(appID), metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return &domain.Container{
				ApplicationID: appID,
				State:         domain.ContainerStateMissing,
			}, nil
		}
		return nil, errors.Wrap(err, "failed to fetch statefulset")
	}
	if isIdleApp(sts) {
		return &domain.Container{
			ApplicationID: appID,
			State:         domain.ContainerStateIdle,
		}, nil
	}

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
	state, msg := getContainerState(list.Items[0].Status)
	return &domain.Container{
		ApplicationID: appID,
		State:         state,
		Message:       msg,
	}, nil
}

func (b *Backend) ListContainers(ctx context.Context) ([]*domain.Container, error) {
	list, err := b.client.CoreV1().Pods(b.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: toSelectorString(allSelector()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pods")
	}
	apps := lo.Map(list.Items, func(pod corev1.Pod, _ int) *domain.Container {
		state, msg := getContainerState(pod.Status)
		return &domain.Container{
			ApplicationID: pod.Labels[appIDLabel],
			State:         state,
			Message:       msg,
		}
	})

	// list statefulsets that enable sablier
	stsList, err := b.client.AppsV1().StatefulSets(b.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: toSelectorString(sablierSelector()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch statefulsets")
	}
	idleApps := lo.FilterMap(stsList.Items, func(sts appv1.StatefulSet, _ int) (*domain.Container, bool) {
		if isIdleApp(&sts) {
			return &domain.Container{
				ApplicationID: sts.Labels[appIDLabel],
				State:         domain.ContainerStateIdle,
			}, true
		}
		return nil, false
	})

	return slices.Concat(apps, idleApps), nil
}

func isIdleApp(sts *appv1.StatefulSet) bool {
	return sts.Labels["sabluer.enable"] == "true" && sts.Status.Replicas == 0
}

func getContainerState(status corev1.PodStatus) (state domain.ContainerState, message string) {
	cs, ok := lo.Find(status.ContainerStatuses, func(cs corev1.ContainerStatus) bool { return cs.Name == podContainerName })
	if !ok {
		return domain.ContainerStateMissing, ""
	}
	if cs.State.Waiting != nil {
		if cs.LastTerminationState.Terminated != nil {
			return domain.ContainerStateRestarting, terminatedMessage(cs.LastTerminationState.Terminated)
		}
		return domain.ContainerStateStarting, waitingMessage(cs.State.Waiting)
	}
	if cs.State.Running != nil {
		return domain.ContainerStateRunning, runningMessage(cs.State.Running)
	}
	if cs.State.Terminated != nil {
		if cs.State.Terminated.ExitCode == 0 {
			return domain.ContainerStateExited, terminatedMessage(cs.State.Terminated)
		} else {
			return domain.ContainerStateErrored, terminatedMessage(cs.State.Terminated)
		}
	}
	return domain.ContainerStateUnknown, "internal error: state unknown"
}

func waitingMessage(state *corev1.ContainerStateWaiting) string {
	msg := state.Reason
	if state.Message != "" {
		msg += ": "
		msg += state.Message
	}
	return msg
}

func runningMessage(state *corev1.ContainerStateRunning) string {
	runningFor := time.Since(state.StartedAt.Time)
	return "Running for " + fmtutil.DurationHuman(runningFor)
}

func terminatedMessage(state *corev1.ContainerStateTerminated) string {
	msg := fmt.Sprintf("Exited with status %d", state.ExitCode)
	if state.Signal != 0 {
		msg += fmt.Sprintf(" (signal %d)", state.Signal)
	}
	if state.Reason != "" {
		msg += ": "
		msg += state.Reason
	}
	if state.Message != "" {
		msg += ": "
		msg += state.Message
	}
	return msg
}
