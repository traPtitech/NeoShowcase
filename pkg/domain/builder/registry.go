package builder

import (
	"context"
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

type RegistryClient interface {
	DeleteImage(ctx context.Context, image, tag string) error
	GetTags(ctx context.Context, image string) ([]string, error)
	GetImageSize(ctx context.Context, image, tag string) (int64, error)
}

func (c *ImageConfig) ImageName(appID string) string {
	return c.Registry.Addr + "/" + c.NamePrefix + appID
}

func (c *ImageConfig) TmpImageName(appID string) string {
	return c.Registry.Addr + "/" + c.TmpNamePrefix + appID
}
