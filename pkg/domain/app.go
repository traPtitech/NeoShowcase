package domain

import (
	"net/url"
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

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

func ExtractNameFromRepositoryURL(repositoryURL string) (string, error) {
	u, err := url.Parse(repositoryURL)
	if err != nil {
		return "", err
	}
	path := u.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ".git")
	return path, nil
}
