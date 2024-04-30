package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/hash"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

const (
	// managedLabel (always "true") indicates the resource is managed by NeoShowcase.
	managedLabel = "ns.trap.jp/managed"
	// appIDLabel indicates the related application ID.
	appIDLabel = "ns.trap.jp/app-id"
	// appRestartAnnotation instructs StatefulSets to restart pods when necessary.
	appRestartAnnotation = "ns.trap.jp/restarted-at"
	// resourceHashAnnotation is hex-encoded 64-bit XXH3 hash of the resource before this annotation is applied.
	resourceHashAnnotation = "ns.trap.jp/hash"
	// fieldManager is the name of this controller.
	fieldManager = "neoshowcase"
)

var _ domain.Backend = (*Backend)(nil)

type Backend struct {
	restConfig        *rest.Config
	client            *kubernetes.Clientset
	traefikClient     *traefikv1alpha1.TraefikV1alpha1Client
	certManagerClient *certmanagerv1.Clientset
	config            Config

	eventSubs   domain.PubSub[*domain.ContainerEvent]
	stopWatcher func()

	reloadLock sync.Mutex
}

func NewK8SBackend(
	restConfig *rest.Config,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikV1alpha1Client,
	certManagerClient *certmanagerv1.Clientset,
	config Config,
) (*Backend, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	b := &Backend{
		restConfig:        restConfig,
		client:            k8sCSet,
		traefikClient:     traefikClient,
		certManagerClient: certManagerClient,
		config:            config,
	}
	return b, nil
}

func (b *Backend) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	b.stopWatcher = cancel
	go retry.Do(ctx, b.eventListener, "pod watcher")
	return nil
}

func (b *Backend) eventListener(ctx context.Context) error {
	podWatcher, err := b.client.CoreV1().Pods(b.config.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			managedLabel: "true",
		}}),
	})
	if err != nil {
		return errors.Wrap(err, "failed to watch pods")
	}
	defer podWatcher.Stop()

	for ev := range podWatcher.ResultChan() {
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
	return nil
}

func (b *Backend) Dispose(_ context.Context) error {
	b.stopWatcher()
	return nil
}

func (b *Backend) AvailableDomains() domain.AvailableDomainSlice {
	return ds.Map(b.config.Domains, (*domainConf).toDomainAD)
}

func (b *Backend) targetAuth(fqdn string) *domainAuthConf {
	for _, dc := range b.config.Domains {
		if dc.Auth.Available && dc.toDomainAD().Match(fqdn) {
			return dc.Auth
		}
	}
	return nil
}

func (b *Backend) AvailablePorts() domain.AvailablePortSlice {
	return ds.Map(b.config.Ports, (*portConf).toDomainAP)
}

func (b *Backend) ListenContainerEvents() (sub <-chan *domain.ContainerEvent, unsub func()) {
	return b.eventSubs.Subscribe()
}

// generalLabelWithoutManagement returns labels that indicates not directly managed by this backend
func (b *Backend) generalLabelWithoutManagement() map[string]string {
	return b.config.labels()
}

func (b *Backend) generalLabel() map[string]string {
	return ds.MergeMap(b.config.labels(), map[string]string{
		managedLabel: "true",
	})
}

func (b *Backend) appLabel(appID string) map[string]string {
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

func deploymentNameWithDiscriminator(appID string, content []byte) string {
	return fmt.Sprintf("nsapp-%s-%s", appID, hash.XXH3Hex(content)[:6])
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
	wildcard := strings.HasPrefix(fqdn, "*.")
	if wildcard {
		fqdn = strings.TrimPrefix(fqdn, "*.")
	}
	fqdn = strings.ReplaceAll(fqdn, "-", "--")
	fqdn = strings.ReplaceAll(fqdn, ".", "-")
	if wildcard {
		return fmt.Sprintf("nsapp-%s-wildcard", fqdn)
	} else {
		return fmt.Sprintf("nsapp-%s", fqdn)
	}
}

func tlsSecretName(fqdn string) string {
	return certificateName(fqdn) + "-tls"
}
