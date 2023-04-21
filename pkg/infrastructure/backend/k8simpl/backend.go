package k8simpl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikcontainous/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	tlsTypeTraefik     = "traefik"
	tlsTypeCertManager = "cert-manager"
)

type authConf = struct {
	Domain   string   `mapstructure:"domain" yaml:"domain"`
	AuthSoft []string `mapstructure:"authSoft" yaml:"authSoft"`
	AuthHard []string `mapstructure:"authHard" yaml:"authHard"`
}

type labelConf = struct {
	Key   string `mapstructure:"key" yaml:"key"`
	Value string `mapstructure:"value" yaml:"value"`
}

type Config struct {
	Middlewares struct {
		Auth []*authConf `mapstructure:"auth" yaml:"auth"`
	} `mapstructure:"middlewares" yaml:"middlewares"`
	SS struct {
		Namespace string `mapstructure:"namespace" yaml:"namespace"`
		Kind      string `mapstructure:"kind" yaml:"kind"`
		Name      string `mapstructure:"name" yaml:"name"`
		Port      int    `mapstructure:"port" yaml:"port"`
	} `mapstructure:"ss" yaml:"ss"`
	Namespace string       `mapstructure:"namespace" yaml:"namespace"`
	Labels    []*labelConf `mapstructure:"labels" yaml:"labels"`
	TLS       struct {
		// cert-manager note: https://doc.traefik.io/traefik/providers/kubernetes-crd/#letsencrypt-support-with-the-custom-resource-definition-provider
		// needs to enable ingress provider in traefik
		Type    string `mapstructure:"type" yaml:"type"`
		Traefik struct {
			CertResolver string `mapstructure:"certResolver" yaml:"certResolver"`
			Wildcard     struct {
				Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
			} `mapstructure:"wildcard" yaml:"wildcard"`
		} `mapstructure:"traefik" yaml:"traefik"`
		CertManager struct {
			Issuer struct {
				Name string `mapstructure:"name" yaml:"name"`
				Kind string `mapstructure:"kind" yaml:"kind"`
			} `mapstructure:"issuer" yaml:"issuer"`
			Wildcard struct {
				Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
			} `mapstructure:"wildcard" yaml:"wildcard"`
		} `mapstructure:"certManager" yaml:"certManager"`
	} `mapstructure:"tls" yaml:"tls"`
	// ImagePullSecret required if registry is private
	ImagePullSecret string `mapstructure:"imagePullSecret" yaml:"imagePullSecret"`
}

func (c *Config) labels() map[string]string {
	return lo.SliceToMap(c.Labels, func(l *labelConf) (string, string) {
		return l.Key, l.Value
	})
}

func (c *Config) Validate() error {
	for _, ac := range c.Middlewares.Auth {
		ad := domain.AvailableDomain{Domain: ac.Domain}
		if err := ad.Validate(); err != nil {
			return errors.Wrapf(err, "invalid domain %s for middleware config", ac.Domain)
		}
	}
	switch c.TLS.Type {
	case tlsTypeTraefik:
		if err := c.TLS.Traefik.Wildcard.Domains.Validate(); err != nil {
			return errors.Wrap(err, "k8s.tls.traefik.wildcard.domains is invalid")
		}
	case tlsTypeCertManager:
		if err := c.TLS.CertManager.Wildcard.Domains.Validate(); err != nil {
			return errors.Wrap(err, "k8s.tls.certManager.wildcard.domains is invalid")
		}
	default:
		return errors.New("k8s.tls.type needs to be one of 'traefik' or 'cert-manager'")
	}
	return nil
}

const (
	appLabel             = "neoshowcase.trap.jp/app"
	appIDLabel           = "neoshowcase.trap.jp/appId"
	appRestartAnnotation = "neoshowcase.trap.jp/startedAt"
	fieldManager         = "neoshowcase"
)

type k8sBackend struct {
	client            *kubernetes.Clientset
	traefikClient     *traefikv1alpha1.TraefikContainousV1alpha1Client
	certManagerClient *certmanagerv1.Clientset
	config            Config
	eventSubs         domain.PubSub[*domain.ContainerEvent]

	podWatcher watch.Interface
	reloadLock sync.Mutex
}

func NewK8SBackend(
	k8sCSet *kubernetes.Clientset,
	traefikClient *traefikv1alpha1.TraefikContainousV1alpha1Client,
	certManagerClient *certmanagerv1.Clientset,
	config Config,
) (domain.Backend, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return &k8sBackend{
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
		b.eventSubs.Publish(&domain.ContainerEvent{ApplicationID: appID})
	}
}

func (b *k8sBackend) Dispose(_ context.Context) error {
	b.podWatcher.Stop()
	return nil
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
