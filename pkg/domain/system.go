package domain

type SSHConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
}

type AdditionalLink struct {
	Name string `mapstructure:"name" yaml:"name"`
	URL  string `mapstructure:"url" yaml:"url"`
}
