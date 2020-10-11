package common

type GRPCConfig struct {
	Port int `mapstructure:"port" yaml:"port"`
	Web  struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"web" yaml:"web"`
}

func (c *GRPCConfig) GetPort() int {
	if c.Port == 0 {
		return 10000
	}
	return c.Port
}
