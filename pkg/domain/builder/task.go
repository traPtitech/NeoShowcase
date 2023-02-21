package builder

type Task struct {
	BuildID       string
	ApplicationID string
	Static        bool
	BuildSource   *BuildSource
	BuildOptions  *BuildOptions
	ImageName     string
}

type BuildSource struct {
	RepositoryUrl string
	Ref           string
}

type BuildOptions struct {
	BaseImageName string
	EntrypointCmd string
	StartupCmd    string
	ArtifactPath  string
}
