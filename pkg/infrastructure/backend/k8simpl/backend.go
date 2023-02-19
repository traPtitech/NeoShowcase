package k8simpl

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
)

const (
	appNamespace                   = "neoshowcase-apps"
	appContainerLabel              = "neoshowcase.trap.jp/app"
	appContainerApplicationIDLabel = "neoshowcase.trap.jp/appId"
	deploymentRestartAnnotation    = "neoshowcase.trap.jp/restartedAt"
)

type k8sBackend struct {
	clientset  *kubernetes.Clientset
	eventbus   domain.Bus
	podWatcher watch.Interface
}

func NewK8SBackend(eventbus domain.Bus, k8sCSet *kubernetes.Clientset) (domain.Backend, error) {
	b := &k8sBackend{
		clientset: k8sCSet,
		eventbus:  eventbus,
	}

	var err error
	b.podWatcher, err = k8sCSet.CoreV1().Pods(appNamespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appContainerLabel: "true",
		}}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to watch pods: %w", err)
	}

	go b.eventListener()

	return b, nil
}

func (b *k8sBackend) eventListener() {
	for ev := range b.podWatcher.ResultChan() {
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
				b.eventbus.Publish(event.ContainerAppStarted, domain.Fields{
					"application_id": p.Labels[appContainerApplicationIDLabel],
				})
			}
		case watch.Deleted:
			b.eventbus.Publish(event.ContainerAppStopped, domain.Fields{
				"application_id": p.Labels[appContainerApplicationIDLabel],
			})
		}
	}
}

func (b *k8sBackend) Dispose(ctx context.Context) error {
	b.podWatcher.Stop()
	return nil
}

func int32Ptr(i int32) *int32                                           { return &i }
func pathTypePtr(pathType networkingv1.PathType) *networkingv1.PathType { return &pathType }

func deploymentName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}
