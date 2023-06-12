package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	// managedLabel (always "true") indicates the resource is managed by NeoShowcase.
	managedLabel = "ns.trap.jp/managed"
	// appIDLabel indicates the related application ID.
	appIDLabel = "ns.trap.jp/app-id"
	// appRestartAnnotation instructs StatefulSets to restart pods when necessary.
	appRestartAnnotation = "ns.trap.jp/restarted-at"
	// fieldManager is the name of this controller.
	fieldManager = "neoshowcase"
)

type k8sBackend struct {
	restConfig        *rest.Config
	client            *kubernetes.Clientset
	traefikClient     *traefikv1alpha1.TraefikV1alpha1Client
	certManagerClient *certmanagerv1.Clientset
	config            Config

	eventSubs domain.PubSub[*domain.ContainerEvent]

	podWatcher watch.Interface
	reloadLock sync.Mutex
}

func NewK8SBackend(
	restConfig *rest.Config,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikV1alpha1Client,
	certManagerClient *certmanagerv1.Clientset,
	config Config,
) (domain.Backend, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	b := &k8sBackend{
		restConfig:        restConfig,
		client:            k8sCSet,
		traefikClient:     traefikClient,
		certManagerClient: certManagerClient,
		config:            config,
	}
	return b, nil
}

func (b *k8sBackend) Start(_ context.Context) error {
	var err error
	b.podWatcher, err = b.client.CoreV1().Pods(b.config.Namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			managedLabel: "true",
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
		b.eventSubs.Publish(&domain.ContainerEvent{ApplicationID: appID})
	}
}

func (b *k8sBackend) Dispose(_ context.Context) error {
	b.podWatcher.Stop()
	return nil
}

func (b *k8sBackend) AvailableDomains() domain.AvailableDomainSlice {
	return ds.Map(b.config.Domains, (*domainConf).toDomainAD)
}

func (b *k8sBackend) targetAuth(fqdn string) *domainAuthConf {
	for _, dc := range b.config.Domains {
		if dc.Auth.Available && dc.toDomainAD().Match(fqdn) {
			return dc.Auth
		}
	}
	return nil
}

func (b *k8sBackend) AvailablePorts() domain.AvailablePortSlice {
	return ds.Map(b.config.Ports, (*portConf).toDomainAP)
}

func (b *k8sBackend) ListenContainerEvents() (sub <-chan *domain.ContainerEvent, unsub func()) {
	return b.eventSubs.Subscribe()
}

// generalLabelWithoutManagement returns labels that indicates not directly managed by this backend
func (b *k8sBackend) generalLabelWithoutManagement() map[string]string {
	return b.config.labels()
}

func (b *k8sBackend) generalLabel() map[string]string {
	return ds.MergeMap(b.config.labels(), map[string]string{
		managedLabel: "true",
	})
}

func (b *k8sBackend) appLabel(appID string) map[string]string {
	return ds.MergeMap(b.config.labels(), map[string]string{
		managedLabel: "true",
		appIDLabel:   appID,
	})
}

func toSelectorString(matchLabels map[string]string) string {
	return metav1.FormatLabelSelector(&metav1.LabelSelector{
		MatchLabels: matchLabels,
	})
}

func allSelector() map[string]string {
	return map[string]string{
		managedLabel: "true",
	}
}

func appSelector(appID string) map[string]string {
	return map[string]string{
		appIDLabel: appID,
	}
}

func deploymentName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

func generatedPodName(appID string) string {
	return deploymentName(appID) + "-0"
}

const podContainerName = "app"

func serviceName(website *domain.Website) string {
	return fmt.Sprintf("nsapp-%s", website.ID)
}

func portServiceName(port *domain.PortPublication) string {
	return fmt.Sprintf("nsapp-port-%s-%d", port.Protocol, port.InternetPort)
}

func stripMiddlewareName(website *domain.Website) string {
	return serviceName(website) + "-strip"
}

func ssHeaderMiddlewareName(ss *domain.StaticSite) string {
	return fmt.Sprintf("nsapp-ss-header-%s", ss.Application.ID)
}

func certificateName(fqdn string) string {
	if strings.HasPrefix(fqdn, "*.") {
		fqdn = strings.TrimPrefix(fqdn, "*.")
		fqdn = strings.ReplaceAll(fqdn, ".", "-")
		return fmt.Sprintf("nsapp-%s-wildcard", fqdn)
	}
	fqdn = strings.ReplaceAll(fqdn, ".", "-")
	return fmt.Sprintf("nsapp-%s", fqdn)
}

func tlsSecretName(fqdn string) string {
	return certificateName(fqdn) + "-tls"
}
