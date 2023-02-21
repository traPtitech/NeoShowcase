package domain

import "github.com/traPtitech/neoshowcase/pkg/domain/builder"

type Application struct {
	ID         string
	Repository Repository
	BranchName string
	BuildType  builder.BuildType
}

type Build struct {
	ID            string
	Status        builder.BuildStatus
	ApplicationID string
}

type Environment struct {
	ID            string
	ApplicationID string
	Key           string
	Value         string
}

type Repository struct {
	ID  string
	URL string
}
