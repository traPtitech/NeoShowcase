package k8simpl

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefik/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
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

const (
	traefikHTTPEntrypoint     = "web"
	traefikHTTPSEntrypoint    = "websecure"
	traefikAuthSoftMiddleware = "ns_auth_soft@file"
	traefikAuthHardMiddleware = "ns_auth_hard@file"
	traefikAuthMiddleware     = "ns_auth@file"
	traefikCertResolver       = "nsresolver@file"
)

type k8sBackend struct {
	client        *kubernetes.Clientset
	traefikClient *traefikv1alpha1.TraefikV1alpha1Client
	eventbus      domain.Bus

	podWatcher watch.Interface
}

func NewK8SBackend(
	eventbus domain.Bus,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikV1alpha1Client,
) (domain.Backend, error) {
	b := &k8sBackend{
		client:        k8sCSet,
		traefikClient: traefikClient,
		eventbus:      eventbus,
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

func (b *k8sBackend) Dispose(_ context.Context) error {
	b.podWatcher.Stop()
	return nil
}

func deploymentName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

func serviceName(fqdn string) string {
	return fmt.Sprintf("nsapp-%s", strings.ReplaceAll(fqdn, ".", "-"))
}
