package builder

type DockerImageRegistryString string
type DockerImageNamePrefixString string

func GetImageName(registry string, prefix string, appID string) string {
	return registry + "/" + prefix + appID
}
