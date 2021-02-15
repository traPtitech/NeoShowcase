package k8simpl

import (
	"context"
	"fmt"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/event"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

const (
	appNamespace                   = "neoshowcase-apps"
	appContainerLabel              = "neoshowcase.trap.jp/app"
	appContainerApplicationIDLabel = "neoshowcase.trap.jp/appId"
	appContainerEnvironmentIDLabel = "neoshowcase.trap.jp/envId"
	deploymentRestartAnnotation    = "neoshowcase.trap.jp/restartedAt"
)

type Manager struct {
	clientset  *kubernetes.Clientset
	bus        *hub.Hub
	podWatcher watch.Interface
}

func NewManager(eventbus *hub.Hub, k8sCSet *kubernetes.Clientset) (*Manager, error) {
	m := &Manager{
		clientset: k8sCSet,
		bus:       eventbus,
	}

	var err error
	m.podWatcher, err = k8sCSet.CoreV1().Pods(appNamespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appContainerLabel: "true",
		}}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to watch pods: %w", err)
	}

	go m.eventListener()

	return m, nil
}

func (m *Manager) eventListener() {
	for ev := range m.podWatcher.ResultChan() {
		p, ok := ev.Object.(*apiv1.Pod)
		if !ok {
			log.Warnf("unexpected type: %v", ev)
			continue
		}
		if p.Labels[appContainerLabel] != "true" {
			continue
		}

		switch ev.Type {
		case watch.Modified:
			if p.Status.Phase == apiv1.PodRunning {
				m.bus.Publish(hub.Message{
					Name: event.ContainerAppStarted,
					Fields: map[string]interface{}{
						"application_id": p.Labels[appContainerApplicationIDLabel],
						"environment_id": p.Labels[appContainerEnvironmentIDLabel],
					},
				})
			}
		case watch.Deleted:
			m.bus.Publish(hub.Message{
				Name: event.ContainerAppStopped,
				Fields: map[string]interface{}{
					"application_id": p.Labels[appContainerApplicationIDLabel],
					"environment_id": p.Labels[appContainerEnvironmentIDLabel],
				},
			})
		}
	}
}

func (m *Manager) Dispose(ctx context.Context) error {
	m.podWatcher.Stop()
	return nil
}
