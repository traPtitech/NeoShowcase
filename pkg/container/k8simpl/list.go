package k8simpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) List(ctx context.Context) (*container.ListResult, error) {
	list, err := m.clientset.CoreV1().Pods(appNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appContainerLabel: "true",
		}}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pods: %w", err)
	}

	var result []container.Container
	for _, item := range list.Items {
		result = append(result, container.Container{
			ApplicationID: item.Labels[appContainerApplicationIDLabel],
			EnvironmentID: item.Labels[appContainerEnvironmentIDLabel],
			State:         getContainerState(item.Status),
		})
	}

	return &container.ListResult{
		Containers: result,
	}, nil
}

func getContainerState(status apiv1.PodStatus) container.State {
	switch status.Phase {
	case apiv1.PodPending:
		return container.StateRestarting
	case apiv1.PodRunning:
		return container.StateRunning
	case apiv1.PodFailed:
		return container.StateStopped
	case apiv1.PodSucceeded:
		return container.StateStopped
	case apiv1.PodUnknown:
		return container.StateOther
	default:
		return container.StateOther
	}
}
