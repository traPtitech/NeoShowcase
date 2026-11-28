package k8srunner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
)

func TestNewRunner_RejectsInvalidConfig(t *testing.T) {
	valid := func() Config {
		var c Config
		c.Namespace = "builds"
		c.Resources.Requests.CPU = "10m"
		c.Resources.Limits.Memory = "1Gi"
		return c
	}
	tests := []struct {
		name       string
		modify     func(c *Config)
		pullPolicy isolatedrunner.PullPolicy
		wantErr    string
	}{
		{name: "valid", modify: func(*Config) {}, pullPolicy: isolatedrunner.PullAlways},
		{name: "no resources", modify: func(c *Config) { *c = Config{Namespace: "builds"} }, pullPolicy: isolatedrunner.PullNever},
		{name: "no namespace", modify: func(c *Config) { c.Namespace = "" }, pullPolicy: isolatedrunner.PullAlways, wantErr: "namespace"},
		{name: "invalid quantity", modify: func(c *Config) { c.Resources.Limits.CPU = "2x" }, pullPolicy: isolatedrunner.PullAlways, wantErr: "limits.cpu"},
		{name: "invalid toleration operator", modify: func(c *Config) { c.Tolerations = []Toleration{{Key: "k", Operator: "Exist"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "Exist"},
		{name: "invalid toleration effect", modify: func(c *Config) { c.Tolerations = []Toleration{{Key: "k", Effect: "NoScheduled"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "NoScheduled"},
		{name: "toleration without operator", modify: func(c *Config) {
			c.Tolerations = []Toleration{{Key: "k", Value: "v", Effect: "NoSchedule"}}
		}, pullPolicy: isolatedrunner.PullAlways},
		{name: "toleration without effect", modify: func(c *Config) { c.Tolerations = []Toleration{{Key: "k", Operator: "Exists"}} }, pullPolicy: isolatedrunner.PullAlways},
		{name: "toleration the API rejects", modify: func(c *Config) { c.Tolerations = []Toleration{{Key: "k", Operator: "Exists", Value: "v"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "value must be empty"},
		{name: "node label with empty key", modify: func(c *Config) { c.NodeSelector = []NodeLabel{{}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "list of key/value"},
		{name: "node label with invalid key", modify: func(c *Config) { c.NodeSelector = []NodeLabel{{Key: "a b", Value: "v"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "a b"},
		{name: "node label with invalid value", modify: func(c *Config) { c.NodeSelector = []NodeLabel{{Key: "k", Value: "a b"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "a b"},
		{name: "node label with empty value", modify: func(c *Config) { c.NodeSelector = []NodeLabel{{Key: "dedicated"}} }, pullPolicy: isolatedrunner.PullAlways},
		{name: "toleration of every taint", modify: func(c *Config) { c.Tolerations = []Toleration{{Operator: "Exists"}} }, pullPolicy: isolatedrunner.PullAlways},
		{name: "toleration with seconds and NoExecute", modify: func(c *Config) {
			c.Tolerations = []Toleration{{Key: "k", Operator: "Exists", Effect: "NoExecute", TolerationSeconds: new(int64(30))}}
		}, pullPolicy: isolatedrunner.PullAlways},
		{name: "toleration with seconds without NoExecute", modify: func(c *Config) {
			c.Tolerations = []Toleration{{Key: "k", Operator: "Exists", Effect: "NoSchedule", TolerationSeconds: new(int64(30))}}
		}, pullPolicy: isolatedrunner.PullAlways, wantErr: "tolerationSeconds"},
		{name: "toleration with invalid key", modify: func(c *Config) { c.Tolerations = []Toleration{{Key: "a b", Operator: "Exists"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "a b"},
		{name: "toleration with invalid value", modify: func(c *Config) { c.Tolerations = []Toleration{{Key: "k", Value: "a b"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "a b"},
		{name: "toleration without operator or key", modify: func(c *Config) { c.Tolerations = []Toleration{{Value: "v"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "key is empty"},
		{name: "toleration Equal without key", modify: func(c *Config) { c.Tolerations = []Toleration{{Operator: "Equal", Value: "v"}} }, pullPolicy: isolatedrunner.PullAlways, wantErr: "key is empty"},
		{name: "invalid pull policy", modify: func(*Config) {}, pullPolicy: "Sometimes", wantErr: "Sometimes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valid()
			tt.modify(&c)
			r, err := NewRunner(nil, nil, Options{Config: c, Hostname: "host", PullPolicy: tt.pullPolicy, ActiveDeadline: time.Minute})
			if tt.wantErr == "" {
				require.NoError(t, err)
				require.NotNil(t, r)
			} else {
				require.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

func TestNewRunner_RejectsActiveDeadlineShorterThanASecond(t *testing.T) {
	for _, d := range []time.Duration{0, 500 * time.Millisecond} {
		_, err := NewRunner(nil, nil, Options{Config: Config{Namespace: "builds"}, Hostname: "host", PullPolicy: isolatedrunner.PullAlways, ActiveDeadline: d})
		require.ErrorContains(t, err, "active deadline")
	}
	_, err := NewRunner(nil, nil, Options{Config: Config{Namespace: "builds"}, Hostname: "host", PullPolicy: isolatedrunner.PullAlways, ActiveDeadline: time.Second})
	require.NoError(t, err)
}
