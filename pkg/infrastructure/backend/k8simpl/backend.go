package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

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
	appNamespace         = "neoshowcase-apps"
	appLabel             = "neoshowcase.trap.jp/app"
	appIDLabel           = "neoshowcase.trap.jp/appId"
	appRestartAnnotation = "neoshowcase.trap.jp/startedAt"
	ssLabel              = "neoshowcase.trap.jp/ss"
	fieldManager         = "neoshowcase"
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

	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	ss        domain.StaticServerConnectivityConfig

	podWatcher watch.Interface
	reloadLock sync.Mutex
}

func NewK8SBackend(
	eventbus domain.Bus,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikV1alpha1Client,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	ss domain.StaticServerConnectivityConfig,
) domain.Backend {
	return &k8sBackend{
		client:        k8sCSet,
		traefikClient: traefikClient,
		eventbus:      eventbus,

		appRepo:   appRepo,
		buildRepo: buildRepo,
		ss:        ss,
	}
}

func (b *k8sBackend) Start(_ context.Context) error {
	var err error
	b.podWatcher, err = b.client.CoreV1().Pods(appNamespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appLabel: "true",
		}}),
	})
	if err != nil {
		return fmt.Errorf("failed to watch pods: %w", err)
	}
	go b.eventListener()
	return nil
}

func (b *k8sBackend) eventListener() {
	for ev := range b.podWatcher.ResultChan() {
		p, ok := ev.Object.(*apiv1.Pod)
		if !ok {
			log.Warnf("unexpected type: %v", ev)
			continue
		}
		if p.Labels[appLabel] != "true" {
			continue
		}

		switch ev.Type {
		case watch.Modified:
			if p.Status.Phase == apiv1.PodRunning {
				b.eventbus.Publish(event.ContainerAppStarted, domain.Fields{
					"application_id": p.Labels[appIDLabel],
				})
			}
		case watch.Deleted:
			b.eventbus.Publish(event.ContainerAppStopped, domain.Fields{
				"application_id": p.Labels[appIDLabel],
			})
		}
	}
}

func (b *k8sBackend) Dispose(_ context.Context) error {
	b.podWatcher.Stop()
	return nil
}

func resourceLabels(appID string) map[string]string {
	return map[string]string{
		appLabel:   "true",
		appIDLabel: appID,
	}
}

func ssResourceLabels(appID string) map[string]string {
	return map[string]string{
		appLabel:   "true",
		appIDLabel: appID,
		ssLabel:    "true",
	}
}

func allSelector() string {
	return metav1.FormatLabelSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			appLabel: "true",
		},
	})
}

func labelSelector(appID string) string {
	return metav1.FormatLabelSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			appIDLabel: appID,
		},
	})
}

func ssLabelSelector() string {
	return metav1.FormatLabelSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			ssLabel: "true",
		},
	})
}

func deploymentName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

func serviceName(website *domain.Website) string {
	s := fmt.Sprintf("nsapp-%s%s",
		strings.ReplaceAll(website.FQDN, ".", "-"),
		strings.ReplaceAll(website.PathPrefix, "/", "-"),
	)
	return strings.TrimSuffix(s, "-")
}

func stripMiddlewareName(website *domain.Website) string {
	return serviceName(website) + "-strip"
}

func ssHeaderMiddlewareName(ss *domain.StaticSite) string {
	return fmt.Sprintf("nsapp-ss-header-%s", ss.Application.ID)
}

func tlsSecretName(fqdn string) string {
	return fmt.Sprintf("nsapp-secret-%s", strings.ReplaceAll(fqdn, ".", "-"))
}
