package dockerimpl

import (
	"math"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

type domainAuthConf = struct {
	Available bool     `mapstructure:"available" yaml:"available"`
	Soft      []string `mapstructure:"soft" yaml:"soft"`
	Hard      []string `mapstructure:"hard" yaml:"hard"`
}

type domainConf struct {
	Domain   string          `mapstructure:"domain" yaml:"domain"`
	Excludes []string        `mapstructure:"excludes" yaml:"excludes"`
	Auth     *domainAuthConf `mapstructure:"auth" yaml:"auth"`
}

func (dc *domainConf) toDomainAD() *domain.AvailableDomain {
	return &domain.AvailableDomain{
		Domain:         dc.Domain,
		ExcludeDomains: dc.Excludes,
		AuthAvailable:  dc.Auth.Available,
	}
}

type portConf struct {
	StartPort int    `mapstructure:"startPort" yaml:"startPort"`
	EndPort   int    `mapstructure:"endPort" yaml:"endPort"`
	Protocol  string `mapstructure:"protocol" yaml:"protocol"`
}

var portProtocolMapper = mapper.MustNewValueMapper(map[string]domain.PortPublicationProtocol{
	"tcp": domain.PortPublicationProtocolTCP,
	"udp": domain.PortPublicationProtocolUDP,
})

func (pc *portConf) toDomainAP() *domain.AvailablePort {
	return &domain.AvailablePort{
		StartPort: pc.StartPort,
		EndPort:   pc.EndPort,
		Protocol:  portProtocolMapper.IntoMust(pc.Protocol),
	}
}

type labelConf = struct {
	Key   string `mapstructure:"key" yaml:"key"`
	Value string `mapstructure:"value" yaml:"value"`
}

type Config struct {
	// ConfDir is the directory to put all traefik configurations in.
	ConfDir string `mapstructure:"confDir" yaml:"confDir"`
	// Domains define available domains to be used by user apps.
	Domains []*domainConf `mapstructure:"domains" yaml:"domains"`
	// Ports define available port-forward ports to be used by user apps.
	Ports []*portConf `mapstructure:"ports" yaml:"ports"`

	// SS defines static-site server endpoint.
	SS struct {
		URL string `mapstructure:"url" yaml:"url"`
	} `mapstructure:"ss" yaml:"ss"`

	// TLS section defines tls setting for user app ingress.
	TLS struct {
		CertResolver string `mapstructure:"certResolver" yaml:"certResolver"`
		Wildcard     struct {
			// Domains define for which (wildcard) domains cert-manager supports configuring DNS records.
			Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
		} `mapstructure:"wildcard" yaml:"wildcard"`
	} `mapstructure:"tls" yaml:"tls"`

	// Network defines which docker network to use with all user apps.
	Network string `mapstructure:"network" yaml:"network"`
	// Labels define common container labels to put to all user apps.
	Labels []*labelConf `mapstructure:"labels" yaml:"labels"`
	// Resources define user app container resource constraints.
	Resources struct {
		CPUs              float64 `mapstructure:"cpus" yaml:"cpus"`
		Memory            int64   `mapstructure:"memory" yaml:"memory"`
		MemorySwap        int64   `mapstructure:"memorySwap" yaml:"memorySwap"`
		MemoryReservation int64   `mapstructure:"memoryReservation" yaml:"memoryReservation"`
	} `mapstructure:"resources" yaml:"resources"`
}

func (c *Config) labels() map[string]string {
	return lo.SliceToMap(c.Labels, func(l *labelConf) (string, string) {
		return l.Key, l.Value
	})
}

func (c *Config) Validate() error {
	for _, dc := range c.Domains {
		if err := dc.toDomainAD().Validate(); err != nil {
			return errors.Wrap(err, "invalid domain config")
		}
	}
	for _, pc := range c.Ports {
		if err := pc.toDomainAP().Validate(); err != nil {
			return errors.Wrap(err, "invalid port config")
		}
	}

	if err := c.TLS.Wildcard.Domains.Validate(); err != nil {
		return errors.Wrap(err, "docker.tls.wildcard.domains is invalid")
	}

	if c.Resources.CPUs < 0 || math.IsNaN(c.Resources.CPUs) || math.IsInf(c.Resources.CPUs, 0) {
		return errors.New("docker.resources.cpus needs to be a positive number")
	}
	if c.Resources.Memory < 0 {
		return errors.New("docker.resources.memory needs to be a positive number")
	}
	if c.Resources.MemorySwap < -1 {
		return errors.New("docker.resources.memorySwap needs to be a positive number, or -1 for unlimited swap")
	}
	if c.Resources.MemoryReservation < 0 {
		return errors.New("docker.resources.memoryReservation needs to be a positive number")
	}

	return nil
}
