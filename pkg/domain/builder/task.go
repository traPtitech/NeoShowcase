package builder

type Task struct {
	BuildID       string
	ApplicationID string
	Static        bool
	BuildSource   *BuildSource
	BuildOptions  *BuildOptions
	ImageName     string
	ImageTag      string
}

type BuildSource struct {
	RepositoryUrl string
	Commit        string
}

type BuildOptions struct {
	BaseImageName  string
	DockerfileName string
	ArtifactPath   string
	BuildCmd       string
	EntrypointCmd  string
}
