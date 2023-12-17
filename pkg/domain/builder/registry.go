package builder

import (
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
)

type RegistryConfig struct {
	Scheme   string `mapstructure:"scheme" yaml:"scheme"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

type ImageConfig struct {
	Registry      RegistryConfig `mapstructure:"registry" yaml:"registry"`
	NamePrefix    string         `mapstructure:"namePrefix" yaml:"namePrefix"`
	TmpNamePrefix string         `mapstructure:"tmpNamePrefix" yaml:"tmpNamePrefix"`
}

// NewRegistry generates a new regclient instance.
func (c *ImageConfig) NewRegistry() *regclient.RegClient {
	var opts []regclient.Opt

	host := config.HostNewName(c.Registry.Scheme + "://" + c.Registry.Addr)
	// RepoAuth should be set to true, because by default regclient internally merges scopes for all repositories
	// it accesses, resulting in a bloating "Authorization" header when accessing large number of repositories at once.
	// also see: https://distribution.github.io/distribution/spec/auth/jwt/
	host.RepoAuth = true
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
