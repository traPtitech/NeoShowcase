package domain

import (
	"net/url"
	"strings"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
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
	StartedAt     time.Time
	FinishedAt    optional.Of[time.Time]
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
