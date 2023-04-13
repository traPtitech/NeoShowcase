package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikcontainous/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	tlsTypeTraefik     = "traefik"
	tlsTypeCertManager = "cert-manager"
)

type Config struct {
	SS struct {
		Namespace string `mapstructure:"namespace" yaml:"namespace"`
		Kind      string `mapstructure:"kind" yaml:"kind"`
		Name      string `mapstructure:"name" yaml:"name"`
		Port      int    `mapstructure:"port" yaml:"port"`
	} `mapstructure:"ss" yaml:"ss"`
	Namespace string            `mapstructure:"namespace" yaml:"namespace"`
	Labels    map[string]string `mapstructure:"labels" yaml:"labels"`
	TLS       struct {
		// cert-manager note: https://doc.traefik.io/traefik/providers/kubernetes-crd/#letsencrypt-support-with-the-custom-resource-definition-provider
		// needs to enable ingress provider in traefik
		Type    string `mapstructure:"type" yaml:"type"`
		Traefik struct {
			CertResolver string `mapstructure:"certResolver" yaml:"certResolver"`
			Wildcard     bool   `mapstructure:"wildcard" yaml:"wildcard"`
		} `mapstructure:"traefik" yaml:"traefik"`
		CertManager struct {
			Issuer struct {
				Name string `mapstructure:"name" yaml:"name"`
				Kind string `mapstructure:"kind" yaml:"kind"`
			} `mapstructure:"issuer" yaml:"issuer"`
			Wildcard bool `mapstructure:"wildcard" yaml:"wildcard"`
		} `mapstructure:"certManager" yaml:"certManager"`
	} `mapstructure:"tls" yaml:"tls"`
	// ImagePullSecret required if registry is private
	ImagePullSecret string `mapstructure:"imagePullSecret" yaml:"imagePullSecret"`
}

const (
	appLabel             = "neoshowcase.trap.jp/app"
	appIDLabel           = "neoshowcase.trap.jp/appId"
	appRestartAnnotation = "neoshowcase.trap.jp/startedAt"
	ssLabel              = "neoshowcase.trap.jp/ss"
	fieldManager         = "neoshowcase"
)

type k8sBackend struct {
	eventbus          domain.Bus
	client            *kubernetes.Clientset
	traefikClient     *traefikv1alpha1.TraefikContainousV1alpha1Client
	certManagerClient *certmanagerv1.Clientset
	config            Config

	podWatcher watch.Interface
	reloadLock sync.Mutex
}

func NewK8SBackend(
	eventbus domain.Bus,
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikContainousV1alpha1Client,
	certManagerClient *certmanagerv1.Clientset,
	config Config,
) (domain.Backend, error) {
	if config.TLS.Type != tlsTypeTraefik && config.TLS.Type != tlsTypeCertManager {
		return nil, errors.New("k8s.tls.type needs to be one of 'traefik' or 'cert-manager'")
	}

	return &k8sBackend{
		eventbus:          eventbus,
		client:            k8sCSet,
		traefikClient:     traefikClient,
		certManagerClient: certManagerClient,
		config:            config,
	}, nil
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

func (b *k8sBackend) appLabel(appID string) map[string]string {
	return ds.MergeMap(b.config.Labels, map[string]string{
		appLabel:   "true",
		appIDLabel: appID,
	})
}

func (b *k8sBackend) ssLabel(appID string) map[string]string {
	return ds.MergeMap(b.config.Labels, map[string]string{
		appLabel:   "true",
		appIDLabel: appID,
		ssLabel:    "true",
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

func ssSelector() map[string]string {
	return map[string]string{
		ssLabel: "true",
	}
}

func deploymentName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

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
	return tlsSecretName(fqdn)
}

func tlsSecretName(fqdn string) string {
	if strings.HasPrefix(fqdn, "*.") {
		fqdn = strings.TrimPrefix(fqdn, "*.")
		fqdn = strings.ReplaceAll(fqdn, ".", "-")
		return fmt.Sprintf("nsapp-wildcard-tls-%s", fqdn)
	}
	fqdn = strings.ReplaceAll(fqdn, ".", "-")
	return fmt.Sprintf("nsapp-tls-%s", fqdn)
}
