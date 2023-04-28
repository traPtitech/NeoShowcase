package k8simpl

type Config struct {
	Namespace     string `mapstructure:"namespace" yaml:"namespace"`
	PodName       string `mapstructure:"podName" yaml:"podName"`
	ContainerName string `mapstructure:"containerName" yaml:"containerName"`
	LocalDir      string `mapstructure:"localDir" yaml:"localDir"`
	RemoteDir     string `mapstructure:"remoteDir" yaml:"remoteDir"`
}
