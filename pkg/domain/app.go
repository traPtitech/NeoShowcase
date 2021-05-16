package domain

import "github.com/traPtitech/neoshowcase/pkg/domain/builder"

type Application struct {
	ID         string
	Repository Repository
}

type Environment struct {
	ID            string
	ApplicationID string
	BranchName    string
	BuildType     builder.BuildType
}

type Repository struct {
	ID        string
	RemoteURL string
}
