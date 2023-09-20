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
