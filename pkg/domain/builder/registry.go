package builder

import (
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
)

type ImageConfig struct {
	Registry struct {
		Scheme   string `mapstructure:"scheme" yaml:"scheme"`
		Addr     string `mapstructure:"addr" yaml:"addr"`
		Username string `mapstructure:"username" yaml:"username"`
		Password string `mapstructure:"password" yaml:"password"`
	} `mapstructure:"registry" yaml:"registry"`
	NamePrefix    string `mapstructure:"namePrefix" yaml:"namePrefix"`
	TmpNamePrefix string `mapstructure:"tmpNamePrefix" yaml:"tmpNamePrefix"`
}

// NewRegistry generates a new regclient instance.
//
// NOTE: should generate a new instance for each image repository access,
// because it internally stores JWT scopes for each repository it accesses.
// Accessing large number of repositories with a single regclient instance
// will result in a bloating "Authorization" header.
func (c *ImageConfig) NewRegistry() *regclient.RegClient {
	var opts []regclient.Opt

	host := config.HostNewName(c.Registry.Scheme + "://" + c.Registry.Addr)
	if c.Registry.Username != "" {
		host.User = c.Registry.Username
	}
	if c.Registry.Password != "" {
		host.Pass = c.Registry.Password
	}
	opts = append(opts, regclient.WithConfigHost(*host))

	return regclient.New(opts...)
}

func (c *ImageConfig) ImageName(appID string) string {
	return c.Registry.Addr + "/" + c.NamePrefix + appID
}

func (c *ImageConfig) TmpImageName(appID string) string {
	return c.Registry.Addr + "/" + c.TmpNamePrefix + appID
}
