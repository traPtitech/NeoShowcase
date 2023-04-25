package k8simpl

import (
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
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

type nodeSelector = struct {
	Key   string `mapstructure:"key" yaml:"key"`
	Value string `mapstructure:"value" yaml:"value"`
}

type toleration = struct {
	Key               string `mapstructure:"key" yaml:"key"`
	Operator          string `mapstructure:"operator" yaml:"operator"`
	Value             string `mapstructure:"value" yaml:"value"`
	Effect            string `mapstructure:"effect" yaml:"effect"`
	TolerationSeconds int    `mapstructure:"tolerationSeconds" yaml:"tolerationSeconds"`
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
	Scheduling      struct {
		NodeSelector []*nodeSelector `mapstructure:"nodeSelector" yaml:"nodeSelector"`
		Tolerations  []*toleration   `mapstructure:"tolerations" yaml:"tolerations"`
	} `mapstructure:"scheduling" yaml:"scheduling"`
}

func (c *Config) labels() map[string]string {
	return lo.SliceToMap(c.Labels, func(l *labelConf) (string, string) {
		return l.Key, l.Value
	})
}

func (c *Config) podSchedulingNodeSelector() map[string]string {
	if len(c.Scheduling.NodeSelector) == 0 {
		return nil
	}
	return lo.SliceToMap(c.Scheduling.NodeSelector, func(ns *nodeSelector) (string, string) {
		return ns.Key, ns.Value
	})
}

var tolerationOperatorMapper = mapper.MustNewValueMapper(map[string]v1.TolerationOperator{
	string(v1.TolerationOpEqual):  v1.TolerationOpEqual,
	string(v1.TolerationOpExists): v1.TolerationOpExists,
})

var taintEffectMapper = mapper.MustNewValueMapper(map[string]v1.TaintEffect{
	string(v1.TaintEffectNoSchedule):       v1.TaintEffectNoSchedule,
	string(v1.TaintEffectNoExecute):        v1.TaintEffectNoExecute,
	string(v1.TaintEffectPreferNoSchedule): v1.TaintEffectPreferNoSchedule,
})

func (c *Config) podSchedulingTolerations() []v1.Toleration {
	if len(c.Scheduling.Tolerations) == 0 {
		return nil
	}
	return lo.Map(c.Scheduling.Tolerations, func(t *toleration, _ int) v1.Toleration {
		return v1.Toleration{
			Key:               t.Key,
			Operator:          tolerationOperatorMapper.IntoMust(t.Operator),
			Value:             t.Value,
			Effect:            taintEffectMapper.IntoMust(t.Effect),
			TolerationSeconds: lo.ToPtr(int64(t.TolerationSeconds)),
		}
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
	for _, t := range c.Scheduling.Tolerations {
		if _, ok := tolerationOperatorMapper.Into(t.Operator); !ok {
			return errors.Errorf("k8s.scheduling.tolerations: unknown toleration operator: %v", t.Operator)
		}
		if _, ok := taintEffectMapper.Into(t.Effect); !ok {
			return errors.Errorf("k8s.scheduling.tolerations: unknown taint effect: %v", t.Effect)
		}
	}
	return nil
}
