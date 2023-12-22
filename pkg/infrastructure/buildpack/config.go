package buildpack

type Config struct {
	Helper struct {
		Address    string `mapstructure:"address" yaml:"address"`
		ListenPort int    `mapstructure:"listenPort" yaml:"listenPort"`
	} `mapstructure:"helper" yaml:"helper"`
	RemoteDir   string `mapstructure:"remoteDir" yaml:"remoteDir"`
	PlatformAPI string `mapstructure:"platformAPI" yaml:"platformAPI"`
}
