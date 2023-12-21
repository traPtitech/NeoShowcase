package domain

import "github.com/traPtitech/neoshowcase/pkg/domain/builder"

type BuilderSystemInfo struct {
	SSHKey      PrivateKey
	ImageConfig builder.ImageConfig
}

type StartBuildRequest struct {
	Repo  *Repository
	App   *Application
	Envs  []*Environment
	Build *Build
}
