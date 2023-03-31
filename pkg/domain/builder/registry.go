package builder

type ImageConfig struct {
	Registry struct {
		Addr     string `mapstructure:"addr" yaml:"addr"`
		Username string `mapstructure:"username" yaml:"username"`
		Password string `mapstructure:"password" yaml:"password"`
	} `mapstructure:"registry" yaml:"registry"`
	NamePrefix string `mapstructure:"namePrefix" yaml:"namePrefix"`
}

func (c *ImageConfig) ImageName(appID string) string {
	return c.Registry.Addr + "/" + c.NamePrefix + appID
}
