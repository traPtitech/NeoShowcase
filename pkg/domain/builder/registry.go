package builder

import (
	"github.com/heroku/docker-registry-client/registry"
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

func (c *ImageConfig) NewRegistry() (*registry.Registry, error) {
	return registry.New(c.Registry.Scheme+"://"+c.Registry.Addr, c.Registry.Username, c.Registry.Password)
}

func (c *ImageConfig) ImageName(appID string) string {
	return c.Registry.Addr + "/" + c.NamePrefix + appID
}

func (c *ImageConfig) TmpImageName(appID string) string {
	return c.Registry.Addr + "/" + c.TmpNamePrefix + appID
}

func (c *ImageConfig) PartialTmpImageName(appID string) string {
	return c.TmpNamePrefix + appID
}
