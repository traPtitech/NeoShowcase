package k8srunner

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func TestPodReached_FailsWhenThePodEndsBeforeTheStepStarts(t *testing.T) {
	// The active deadline of the pod passes while the receive step is still being prepared.
	pod := &v1.Pod{Status: v1.PodStatus{
		Phase:  v1.PodFailed,
		Reason: "DeadlineExceeded",
		InitContainerStatuses: []v1.ContainerStatus{{
			Name:  "receive",
			State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "PodInitializing"}},
		}},
	}}
	_, err := podReached(pod, stepStarted("detect"))
	require.ErrorContains(t, err, "DeadlineExceeded")
}

func TestPodReached_WaitsWhileThePodMayProceed(t *testing.T) {
	pod := &v1.Pod{Status: v1.PodStatus{
		Phase: v1.PodPending,
		InitContainerStatuses: []v1.ContainerStatus{{
			Name:  "receive",
			State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "PodInitializing"}},
		}},
	}}
	done, err := podReached(pod, stepStarted("detect"))
	require.NoError(t, err)
	require.False(t, done)
}

func TestPodReached_DoneOnceTheConditionHolds(t *testing.T) {
	// A step that has exited counts as started, even after the pod ended.
	pod := &v1.Pod{Status: v1.PodStatus{
		Phase: v1.PodFailed,
		ContainerStatuses: []v1.ContainerStatus{{
			Name:  "detect",
			State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 1}},
		}},
	}}
	done, err := podReached(pod, stepStarted("detect"))
	require.NoError(t, err)
	require.True(t, done)
}
