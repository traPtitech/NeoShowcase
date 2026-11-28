package k8srunner

import (
	"slices"
	"strings"

	"github.com/samber/oops"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// NodeLabel is an entry of the node selector of build pods.
// A node selector is a list rather than a map in the config, since the config loader splits map keys at dots, which label keys contain.
type NodeLabel struct {
	Key   string `mapstructure:"key" yaml:"key"`
	Value string `mapstructure:"value" yaml:"value"`
}

// Toleration is a toleration of build pods. It has the fields of the same names in v1.Toleration, which give an empty Operator and Effect their meanings.
type Toleration struct {
	Key               string `mapstructure:"key" yaml:"key"`
	Operator          string `mapstructure:"operator" yaml:"operator"`
	Value             string `mapstructure:"value" yaml:"value"`
	Effect            string `mapstructure:"effect" yaml:"effect"`
	TolerationSeconds *int64 `mapstructure:"tolerationSeconds" yaml:"tolerationSeconds"`
}

var (
	tolerationOperators = []v1.TolerationOperator{"", v1.TolerationOpEqual, v1.TolerationOpExists}
	taintEffects        = []v1.TaintEffect{"", v1.TaintEffectNoSchedule, v1.TaintEffectPreferNoSchedule, v1.TaintEffectNoExecute}
)

func validateNodeSelector(labels []NodeLabel) error {
	for i, l := range labels {
		if l.Key == "" {
			return oops.Errorf("nodeSelector[%d]: key is empty; write nodeSelector as a list of key/value, not as a map", i)
		}
		if msgs := validation.IsQualifiedName(l.Key); len(msgs) > 0 {
			return oops.Errorf("nodeSelector[%d]: invalid key %q: %s", i, l.Key, strings.Join(msgs, "; "))
		}
		if msgs := validation.IsValidLabelValue(l.Value); len(msgs) > 0 {
			return oops.Errorf("nodeSelector[%d]: invalid value %q: %s", i, l.Value, strings.Join(msgs, "; "))
		}
	}
	return nil
}

func validateTolerations(tolerations []Toleration) error {
	for i, t := range tolerations {
		if err := validateToleration(t); err != nil {
			return oops.Wrapf(err, "tolerations[%d]", i)
		}
	}
	return nil
}

func validateToleration(t Toleration) error {
	operator := v1.TolerationOperator(t.Operator)
	if !slices.Contains(tolerationOperators, operator) {
		return oops.Errorf("unknown operator %q", t.Operator)
	}
	if !slices.Contains(taintEffects, v1.TaintEffect(t.Effect)) {
		return oops.Errorf("unknown effect %q", t.Effect)
	}
	if t.Key == "" {
		if operator != v1.TolerationOpExists {
			return oops.New("key is empty, which only operator Exists allows")
		}
	} else if msgs := validation.IsQualifiedName(t.Key); len(msgs) > 0 {
		return oops.Errorf("invalid key %q: %s", t.Key, strings.Join(msgs, "; "))
	}
	if operator == v1.TolerationOpExists {
		if t.Value != "" {
			return oops.Errorf("value must be empty with operator Exists, but is %q", t.Value)
		}
	} else if msgs := validation.IsValidLabelValue(t.Value); len(msgs) > 0 {
		return oops.Errorf("invalid value %q: %s", t.Value, strings.Join(msgs, "; "))
	}
	if t.TolerationSeconds != nil && v1.TaintEffect(t.Effect) != v1.TaintEffectNoExecute {
		return oops.Errorf("tolerationSeconds is allowed only with effect NoExecute, but effect is %q", t.Effect)
	}
	return nil
}

func podNodeSelector(labels []NodeLabel) map[string]string {
	if len(labels) == 0 {
		return nil
	}
	m := make(map[string]string, len(labels))
	for _, l := range labels {
		m[l.Key] = l.Value
	}
	return m
}

func podTolerations(tolerations []Toleration) []v1.Toleration {
	var list []v1.Toleration
	for _, t := range tolerations {
		list = append(list, v1.Toleration{
			Key:               t.Key,
			Operator:          v1.TolerationOperator(t.Operator),
			Value:             t.Value,
			Effect:            v1.TaintEffect(t.Effect),
			TolerationSeconds: t.TolerationSeconds,
		})
	}
	return list
}
