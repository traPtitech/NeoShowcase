package k8simpl

import (
	"fmt"
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/hash"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

const (
	routingTypeTraefik = "traefik"

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
	// Domains define available domains to be used by user apps.
	Domains []*domainConf `mapstructure:"domains" yaml:"domains"`
	// Ports define available port-forward ports to be used by user apps.
	Ports []*portConf `mapstructure:"ports" yaml:"ports"`

	// SS defines static-site server endpoint.
	SS struct {
		Namespace string `mapstructure:"namespace" yaml:"namespace"`
		Kind      string `mapstructure:"kind" yaml:"kind"`
		Name      string `mapstructure:"name" yaml:"name"`
		Port      int    `mapstructure:"port" yaml:"port"`
		Scheme    string `mapstructure:"scheme" yaml:"scheme"`
	} `mapstructure:"ss" yaml:"ss"`

	// Routing section defines ingress controller settings.
	Routing struct {
		// Type defines which ingress controller to use.
		// Possible values:
		// 	"traefik": Uses traefik ingress controller.
		Type    string `mapstructure:"type" yaml:"type"`
		Traefik struct {
			// PriorityOffset defines HTTP routes' priority offset for user apps.
			// This is optionally used to optimize routing performance.
			PriorityOffset int `mapstructure:"priorityOffset" yaml:"priorityOffset"`
		} `mapstructure:"traefik" yaml:"traefik"`
	} `mapstructure:"routing" yaml:"routing"`
	// Service section defines Service (L4) routing settings.
	Service struct {
		// IPFamilies defines ipFamilies field for the service objects.
		// Allowed values: IPv4, IPv6
		IPFamilies []v1.IPFamily `mapstructure:"ipFamilies" yaml:"ipFamilies"`
		// IPFamilyPolicy defines ipFamilyPolicy field for the service objects.
		// Allowed values: "", "SingleStack", "PreferDualStack", "RequireDualStack"
		IPFamilyPolicy v1.IPFamilyPolicy `mapstructure:"ipFamilyPolicy" yaml:"ipFamilyPolicy"`
	} `mapstructure:"service" yaml:"service"`

	// Middleware section defines middleware settings.
	Middleware struct {
		// Sablier (https://github.com/acouvreur/sablier) starts user apps on demand and shuts them down after a certain time.
		Sablier struct {
			Enable     bool   `mapstructure:"enable" yaml:"enable"`
			SablierURL string `mapstructure:"url" yaml:"url"`
			// SessionDuration defines how long the session should last.
			//
			// Example: "10m"
			SessionDuration string `mapstructure:"sessionDuration" yaml:"sessionDuration"`
		} `mapstructure:"sablier" yaml:"sablier"`
	}
	// TLS section defines tls setting for user app ingress.
	TLS struct {
		// Type defines which provider is responsible for obtaining http certificates.
		// Possible values:
		// 	"traefik: Uses traefik internal Lets Encrypt resolver.
		// 	"cert-manager": Delegates to cert-resolver with its own CRD.
		//
		// NOTE: If using multiple instances of traefik, the traefik internal Lets Encrypt resolver is not supported.
		// https://doc.traefik.io/traefik/providers/kubernetes-crd/#letsencrypt-support-with-the-custom-resource-definition-provider
		Type string `mapstructure:"type" yaml:"type"`

		// Traefik section defines options for type "traefik".
		Traefik struct {
			CertResolver string `mapstructure:"certResolver" yaml:"certResolver"`
			Wildcard     struct {
				Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
			} `mapstructure:"wildcard" yaml:"wildcard"`
		} `mapstructure:"traefik" yaml:"traefik"`

		// CertManager section defines options for type "cert-manager".
		CertManager struct {
			// Issuer defines cert-manager Issuer object reference to be used with CRD.
			Issuer struct {
				Name string `mapstructure:"name" yaml:"name"`
				Kind string `mapstructure:"kind" yaml:"kind"`
			} `mapstructure:"issuer" yaml:"issuer"`
			Wildcard struct {
				// Domains define for which (wildcard) domains cert-manager supports configuring DNS records.
				Domains domain.WildcardDomains `mapstructure:"domains" yaml:"domains"`
			} `mapstructure:"wildcard" yaml:"wildcard"`
		} `mapstructure:"certManager" yaml:"certManager"`
	} `mapstructure:"tls" yaml:"tls"`

	// Namespace defines in which namespace to deploy user apps.
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
	// ImagePullSecret defines secret name to pull user app images with. Required if registry is private.
	ImagePullSecret string `mapstructure:"imagePullSecret" yaml:"imagePullSecret"`
	// Labels define common labels put to all NeoShowcase-managed Kubernetes objects.
	Labels []*labelConf `mapstructure:"labels" yaml:"labels"`
	// Scheduling defines user app pod scheduling constraints.
	Scheduling struct {
		NodeSelector      []*nodeSelector     `mapstructure:"nodeSelector" yaml:"nodeSelector"`
		Tolerations       []*toleration       `mapstructure:"tolerations" yaml:"tolerations"`
		SpreadConstraints []*spreadConstraint `mapstructure:"spreadConstraints" yaml:"spreadConstraints"`
		ForceHosts        []string            `mapstructure:"forceHosts" yaml:"forceHosts"`
	} `mapstructure:"scheduling" yaml:"scheduling"`
	// Resources define user app pod resource constraints.
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

func (c *Config) serviceIPFamilies() []v1.IPFamily {
	return c.Service.IPFamilies
}

func (c *Config) serviceIPFamilyPolicy() *v1.IPFamilyPolicy {
	if c.Service.IPFamilyPolicy == "" {
		return nil
	}
	return lo.ToPtr(c.Service.IPFamilyPolicy)
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

	switch c.Routing.Type {
	case routingTypeTraefik:
		// Nothing to validate as of now
	default:
		return errors.New(fmt.Sprintf("k8s.routing.type is invalid: %s", c.Routing.Type))
	}

	for _, family := range c.Service.IPFamilies {
		if !lo.Contains([]v1.IPFamily{v1.IPv4Protocol, v1.IPv6Protocol}, family) {
			return errors.New(fmt.Sprintf("invalid IPFamily %s", family))
		}
	}
	if !lo.Contains([]v1.IPFamilyPolicy{
		// Allow empty value
		"",
		v1.IPFamilyPolicySingleStack,
		v1.IPFamilyPolicyPreferDualStack,
		v1.IPFamilyPolicyRequireDualStack,
	}, c.Service.IPFamilyPolicy) {
		return errors.New(fmt.Sprintf("invalid IPFamily policy: %s", c.Service.IPFamilyPolicy))
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
