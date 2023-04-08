package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikcontainous/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
)

type Config struct {
	SS struct {
		Namespace string `mapstructure:"namespace" yaml:"namespace"`
		Kind      string `mapstructure:"kind" yaml:"kind"`
		Name      string `mapstructure:"name" yaml:"name"`
		Port      int    `mapstructure:"port" yaml:"port"`
	} `mapstructure:"ss" yaml:"ss"`
}

type (
	m map[string]any
)

const (
	appNamespace         = "neoshowcase-apps"
	appLabel             = "neoshowcase.trap.jp/app"
	appIDLabel           = "neoshowcase.trap.jp/appId"
	appRestartAnnotation = "neoshowcase.trap.jp/startedAt"
	ssLabel              = "neoshowcase.trap.jp/ss"
	fieldManager         = "neoshowcase"
)

type k8sBackend struct {
	eventbus      domain.Bus
	client        *kubernetes.Clientset
	traefikClient *traefikv1alpha1.TraefikContainousV1alpha1Client
	config        Config

	podWatcher watch.Interface
	reloadLock sync.Mutex
}

func NewK8SBackend(
	eventbus domain.Bus,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikContainousV1alpha1Client,
	config Config,
) domain.Backend {
	return &k8sBackend{
		client:        k8sCSet,
		traefikClient: traefikClient,
		eventbus:      eventbus,
		config:        config,
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
		return errors.Wrap(err, "failed to watch pods")
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

		appID, ok := p.Labels[appIDLabel]
		if !ok {
			continue
		}
		b.eventbus.Publish(event.AppContainerUpdated, domain.Fields{
			"application_id": appID,
		})
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
