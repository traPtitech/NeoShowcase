package k8simpl

import (
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/hash"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

const (
	tlsTypeTraefik     = "traefik"
	tlsTypeCertManager = "cert-manager"

	hostnameNodeSelectorLabel = "kubernetes.io/hostname"
)

type middleware struct {
	Name      string `mapstructure:"name" yaml:"name"`
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
}

func (mw *middleware) toRef() traefikv1alpha1.MiddlewareRef {
	return traefikv1alpha1.MiddlewareRef{
		Name:      mw.Name,
		Namespace: mw.Namespace,
	}
}

type domainAuthConf = struct {
	Available bool         `mapstructure:"available" yaml:"available"`
	Soft      []middleware `mapstructure:"soft" yaml:"soft"`
	Hard      []middleware `mapstructure:"hard" yaml:"hard"`
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

type labelSelector = struct {
	MatchExpressions []*labelExpression `mapstructure:"matchExpressions" yaml:"matchExpressions"`
	MatchLabels      []*labelConf       `mapstructure:"matchLabels" yaml:"matchLabels"`
}

func toLabelSelector(l *labelSelector) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: lo.SliceToMap(l.MatchLabels, func(l *labelConf) (string, string) {
			return l.Key, l.Value
		}),
		MatchExpressions: ds.Map(l.MatchExpressions, toLabelExpression),
	}
}

type labelExpression = struct {
	Key      string   `mapstructure:"key" yaml:"key"`
	Operator string   `mapstructure:"operator" yaml:"operator"`
	Values   []string `mapstructure:"values" yaml:"values"`
}

var labelSelectorOperatorMapper = mapper.MustNewValueMapper(map[string]metav1.LabelSelectorOperator{
	string(metav1.LabelSelectorOpIn):           metav1.LabelSelectorOpIn,
	string(metav1.LabelSelectorOpNotIn):        metav1.LabelSelectorOpNotIn,
	string(metav1.LabelSelectorOpExists):       metav1.LabelSelectorOpExists,
	string(metav1.LabelSelectorOpDoesNotExist): metav1.LabelSelectorOpDoesNotExist,
})

func toLabelExpression(l *labelExpression) metav1.LabelSelectorRequirement {
	return metav1.LabelSelectorRequirement{
		Key:      l.Key,
		Operator: labelSelectorOperatorMapper.IntoMust(l.Operator),
		Values:   l.Values,
	}
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
	TolerationSeconds *int64 `mapstructure:"tolerationSeconds" yaml:"tolerationSeconds"`
}

var nodeInclusionPolicyMapper = mapper.MustNewValueMapper(map[string]v1.NodeInclusionPolicy{
	string(v1.NodeInclusionPolicyIgnore): v1.NodeInclusionPolicyIgnore,
	string(v1.NodeInclusionPolicyHonor):  v1.NodeInclusionPolicyHonor,
})

func nodeInclusionPolicyMap(s string) *v1.NodeInclusionPolicy {
	if s == "" {
		return nil
	}
	policy := nodeInclusionPolicyMapper.IntoMust(s)
	return &policy
}

type spreadConstraint = struct {
	MaxSkew            int32          `mapstructure:"maxSkew" yaml:"maxSkew"`
	MinDomains         *int32         `mapstructure:"minDomains" yaml:"minDomains"`
	TopologyKey        string         `mapstructure:"topologyKey" yaml:"topologyKey"`
	WhenUnsatisfiable  string         `mapstructure:"whenUnsatisfiable" yaml:"whenUnsatisfiable"`
	LabelSelector      *labelSelector `mapstructure:"labelSelector" yaml:"labelSelector"`
	MatchLabelKeys     []string       `mapstructure:"matchLabelKeys" yaml:"matchLabelKeys"`
	NodeAffinityPolicy string         `mapstructure:"nodeAffinityPolicy" yaml:"nodeAffinityPolicy"`
	NodeTaintsPolicy   string         `mapstructure:"nodeTaintsPolicy" yaml:"nodeTaintsPolicy"`
}

type Config struct {
	Domains []*domainConf `mapstructure:"domains" yaml:"domains"`
	Ports   []*portConf   `mapstructure:"ports" yaml:"ports"`
	SS      struct {
		Namespace string `mapstructure:"namespace" yaml:"namespace"`
		Kind      string `mapstructure:"kind" yaml:"kind"`
		Name      string `mapstructure:"name" yaml:"name"`
		Port      int    `mapstructure:"port" yaml:"port"`
		Scheme    string `mapstructure:"scheme" yaml:"scheme"`
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
		NodeSelector      []*nodeSelector     `mapstructure:"nodeSelector" yaml:"nodeSelector"`
		Tolerations       []*toleration       `mapstructure:"tolerations" yaml:"tolerations"`
		SpreadConstraints []*spreadConstraint `mapstructure:"spreadConstraints" yaml:"spreadConstraints"`
		ForceHosts        []string            `mapstructure:"forceHosts" yaml:"forceHosts"`
	} `mapstructure:"scheduling" yaml:"scheduling"`
	Resources struct {
		Requests struct {
			CPU    string `mapstructure:"cpu" yaml:"cpu"`
			Memory string `mapstructure:"memory" yaml:"memory"`
		} `mapstructure:"requests" yaml:"requests"`
		Limits struct {
			CPU    string `mapstructure:"cpu" yaml:"cpu"`
			Memory string `mapstructure:"memory" yaml:"memory"`
		} `mapstructure:"limits" yaml:"limits"`
	} `mapstructure:"resources" yaml:"resources"`
}

func (c *Config) labels() map[string]string {
	return lo.SliceToMap(c.Labels, func(l *labelConf) (string, string) {
		return l.Key, l.Value
	})
}

func (c *Config) podSchedulingNodeSelector(appID string) map[string]string {
	ns := lo.SliceToMap(c.Scheduling.NodeSelector, func(ns *nodeSelector) (string, string) {
		return ns.Key, ns.Value
	})
	host := c.selectNode(appID)
	return ds.MergeMap(ns, host)
}

func (c *Config) selectNode(appID string) map[string]string {
	if len(c.Scheduling.ForceHosts) == 0 {
		return nil
	}
	// NOTE: XXH3Hex always returns a 64-bit hex string
	hex, _ := strconv.ParseUint(hash.XXH3Hex([]byte(appID)), 16, 64)
	host := c.Scheduling.ForceHosts[hex%uint64(len(c.Scheduling.ForceHosts))]
	return map[string]string{hostnameNodeSelectorLabel: host}
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

var unsatisfiableMapper = mapper.MustNewValueMapper(map[string]v1.UnsatisfiableConstraintAction{
	string(v1.DoNotSchedule):  v1.DoNotSchedule,
	string(v1.ScheduleAnyway): v1.ScheduleAnyway,
})

func (c *Config) podSchedulingTolerations() []v1.Toleration {
	return ds.Map(c.Scheduling.Tolerations, func(t *toleration) v1.Toleration {
		return v1.Toleration{
			Key:               t.Key,
			Operator:          tolerationOperatorMapper.IntoMust(t.Operator),
			Value:             t.Value,
			Effect:            taintEffectMapper.IntoMust(t.Effect),
			TolerationSeconds: t.TolerationSeconds,
		}
	})
}

func (c *Config) podSpreadConstraints() []v1.TopologySpreadConstraint {
	return ds.Map(c.Scheduling.SpreadConstraints, func(c *spreadConstraint) v1.TopologySpreadConstraint {
		return v1.TopologySpreadConstraint{
			MaxSkew:            c.MaxSkew,
			TopologyKey:        c.TopologyKey,
			WhenUnsatisfiable:  unsatisfiableMapper.IntoMust(c.WhenUnsatisfiable),
			LabelSelector:      toLabelSelector(c.LabelSelector),
			MinDomains:         c.MinDomains,
			NodeAffinityPolicy: nodeInclusionPolicyMap(c.NodeAffinityPolicy),
			NodeTaintsPolicy:   nodeInclusionPolicyMap(c.NodeTaintsPolicy),
			MatchLabelKeys:     c.MatchLabelKeys,
		}
	})
}

func (c *Config) resourceRequirements() v1.ResourceRequirements {
	var r v1.ResourceRequirements
	if c.Resources.Requests.CPU != "" {
		ds.AppendMap(&r.Requests, v1.ResourceCPU, resource.MustParse(c.Resources.Requests.CPU))
	}
	if c.Resources.Requests.Memory != "" {
		ds.AppendMap(&r.Requests, v1.ResourceMemory, resource.MustParse(c.Resources.Requests.Memory))
	}
	if c.Resources.Limits.CPU != "" {
		ds.AppendMap(&r.Limits, v1.ResourceCPU, resource.MustParse(c.Resources.Limits.CPU))
	}
	if c.Resources.Limits.Memory != "" {
		ds.AppendMap(&r.Limits, v1.ResourceMemory, resource.MustParse(c.Resources.Limits.Memory))
	}
	return r
}

func validateResourceQuantity(s string) error {
	_, err := resource.ParseQuantity(s)
	return err
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

	if c.Resources.Requests.CPU != "" {
		if err := validateResourceQuantity(c.Resources.Requests.CPU); err != nil {
			return errors.Wrap(err, "k8s.resources.requests.cpu: invalid quantity")
		}
	}
	if c.Resources.Requests.Memory != "" {
		if err := validateResourceQuantity(c.Resources.Requests.Memory); err != nil {
			return errors.Wrap(err, "k8s.resources.requests.memory: invalid quantity")
		}
	}
	if c.Resources.Limits.CPU != "" {
		if err := validateResourceQuantity(c.Resources.Limits.CPU); err != nil {
			return errors.Wrap(err, "k8s.resources.limits.cpu: invalid quantity")
		}
	}
	if c.Resources.Limits.Memory != "" {
		if err := validateResourceQuantity(c.Resources.Limits.Memory); err != nil {
			return errors.Wrap(err, "k8s.resources.limits.memory: invalid quantity")
		}
	}

	return nil
}
