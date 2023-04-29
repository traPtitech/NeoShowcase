package dockerimpl

type Config struct {
	ContainerName string `mapstructure:"containerName" yaml:"containerName"`
	RemoteDir     string `mapstructure:"remoteDir" yaml:"remoteDir"`
	User          string `mapstructure:"user" yaml:"user"`
	Group         string `mapstructure:"group" yaml:"group"`
	PlatformAPI   string `mapstructure:"platformAPI" yaml:"platformAPI"`
}
