package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
	appLabel             = "neoshowcase.trap.jp/app"
	appIDLabel           = "neoshowcase.trap.jp/appId"
	appRestartAnnotation = "neoshowcase.trap.jp/startedAt"
	fieldManager         = "neoshowcase"
)

type k8sBackend struct {
	restConfig        *rest.Config
	client            *kubernetes.Clientset
	traefikClient     *traefikv1alpha1.TraefikV1alpha1Client
	certManagerClient *certmanagerv1.Clientset
	config            Config
	appRepo           domain.ApplicationRepository
	userRepo          domain.UserRepository

	eventSubs domain.PubSub[*domain.ContainerEvent]
	sshServer *sshServer

	podWatcher watch.Interface
	reloadLock sync.Mutex
}

func NewK8SBackend(
	restConfig *rest.Config,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikV1alpha1Client,
	certManagerClient *certmanagerv1.Clientset,
	config Config,
	key *ssh.PublicKeys,
	appRepo domain.ApplicationRepository,
	userRepo domain.UserRepository,
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
		appRepo:           appRepo,
		userRepo:          userRepo,
	}
	b.sshServer = newSSHServer(b, config.SSH, key)
	return b, nil
}

func (b *k8sBackend) Start(_ context.Context) error {
	var err error
	b.podWatcher, err = b.client.CoreV1().Pods(b.config.Namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: map[string]string{
			appLabel: "true",
		}}),
	})
	if err != nil {
		return errors.Wrap(err, "failed to watch pods")
	}
	go b.eventListener()

	b.sshServer.Start()

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
	return b.sshServer.Close()
}

func (b *k8sBackend) AuthAllowed(fqdn string) bool {
	for _, ac := range b.config.Middlewares.Auth {
		if domain.MatchDomain(ac.Domain, fqdn) {
			return true
		}
	}
	return false
}

func (b *k8sBackend) targetAuth(fqdn string) *authConf {
	for _, ac := range b.config.Middlewares.Auth {
		if domain.MatchDomain(ac.Domain, fqdn) {
			return ac
		}
	}
	return nil
}

func (b *k8sBackend) ListenContainerEvents() (sub <-chan *domain.ContainerEvent, unsub func()) {
	return b.eventSubs.Subscribe()
}

func (b *k8sBackend) generalLabel() map[string]string {
	return ds.MergeMap(b.config.labels(), map[string]string{
		appLabel: "true",
	})
}

func (b *k8sBackend) appLabel(appID string) map[string]string {
	return ds.MergeMap(b.config.labels(), map[string]string{
		appLabel:   "true",
		appIDLabel: appID,
	})
}

func toSelectorString(matchLabels map[string]string) string {
	return metav1.FormatLabelSelector(&metav1.LabelSelector{
		MatchLabels: matchLabels,
	})
}

func allSelector() map[string]string {
	return map[string]string{
		appLabel: "true",
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
