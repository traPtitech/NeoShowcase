package dockerimpl

import (
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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
	ConfDir     string `mapstructure:"confDir" yaml:"confDir"`
	Middlewares struct {
		Auth []*authConf `mapstructure:"auth" yaml:"auth"`
	} `mapstructure:"middlewares" yaml:"middlewares"`
	SS struct {
		URL string `mapstructure:"url" yaml:"url"`
	} `mapstructure:"ss" yaml:"ss"`
	Network string       `mapstructure:"network" yaml:"network"`
	Labels  []*labelConf `mapstructure:"labels" yaml:"labels"`
	TLS     struct {
		CertResolver string `mapstructure:"certResolver" yaml:"certResolver"`
		Wildcard     struct {
			Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
		} `mapstructure:"wildcard" yaml:"wildcard"`
	} `mapstructure:"tls" yaml:"tls"`
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
	if err := c.TLS.Wildcard.Domains.Validate(); err != nil {
		return errors.Wrap(err, "docker.tls.wildcard.domains is invalid")
	}
	return nil
}
