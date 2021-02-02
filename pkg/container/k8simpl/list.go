package k8simpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) List(ctx context.Context) (*container.ListResult, error) {
	list, err := m.clientset.AppsV1().Deployments(appNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appContainerLabel: "true",
		}}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployments: %w", err)
	}

	var result []container.Container
	for _, item := range list.Items {
		result = append(result, container.Container{
			ApplicationID: item.Labels[appContainerApplicationIDLabel],
			EnvironmentID: item.Labels[appContainerEnvironmentIDLabel],
			State:         item.Status.String(), // TODO
		})
	}

	return &container.ListResult{
		Containers: result,
	}, nil
}
