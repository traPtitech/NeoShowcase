package k8simpl

import (
	"context"
	"fmt"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
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
	clientset *kubernetes.Clientset

	deploymentWatcher watch.Interface
}

func NewManager(eventbus *hub.Hub, k8sCSet *kubernetes.Clientset) (*Manager, error) {
	m := &Manager{
		clientset: k8sCSet,
	}

	var err error
	m.deploymentWatcher, err = k8sCSet.AppsV1().Deployments(appNamespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appContainerLabel: "true",
		}}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to watch deployments: %w", err)
	}

	go m.eventListener()

	return m, nil
}

func (m *Manager) eventListener() {
	for ev := range m.deploymentWatcher.ResultChan() {
		log.Debug(ev)
	}
}

func (m *Manager) Dispose(ctx context.Context) error {
	m.deploymentWatcher.Stop()
	return nil
}
