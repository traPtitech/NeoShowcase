package k8simpl

import (
	"context"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/fmtutil"
)

func (b *Backend) GetContainer(ctx context.Context, appID string) (*domain.Container, error) {
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

	result := ds.Map(list.Items, func(pod v1.Pod) *domain.Container {
		state, msg := getContainerState(pod.Status)
		return &domain.Container{
			ApplicationID: pod.Labels[appIDLabel],
			State:         state,
			Message:       msg,
		}
	})
	return result, nil
}

func getContainerState(status v1.PodStatus) (state domain.ContainerState, message string) {
	cs, ok := lo.Find(status.ContainerStatuses, func(cs v1.ContainerStatus) bool { return cs.Name == podContainerName })
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

func waitingMessage(state *v1.ContainerStateWaiting) string {
	msg := state.Reason
	if state.Message != "" {
		msg += ": "
		msg += state.Message
	}
	return msg
}

func runningMessage(state *v1.ContainerStateRunning) string {
	runningFor := time.Since(state.StartedAt.Time)
	return "Running for " + fmtutil.DurationHuman(runningFor)
}

func terminatedMessage(state *v1.ContainerStateTerminated) string {
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
