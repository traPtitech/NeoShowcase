package builder

import "github.com/volatiletech/null/v8"

type Task struct {
	BuildID      string
	BranchID     null.String
	Static       bool
	BuildSource  *BuildSource
	BuildOptions *BuildOptions
	ImageName    string
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
