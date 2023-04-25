package k8simpl

import (
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

const (
	tlsTypeTraefik     = "traefik"
	tlsTypeCertManager = "cert-manager"
)

type authConf = struct {
	Domain string   `mapstructure:"domain" yaml:"domain"`
	Soft   []string `mapstructure:"soft" yaml:"soft"`
	Hard   []string `mapstructure:"hard" yaml:"hard"`
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
