package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ApplicationState int

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateStarting
	ApplicationStateRunning
	ApplicationStateErrored
)

func (s ApplicationState) String() string {
	switch s {
	case ApplicationStateIdle:
		return "IDLE"
	case ApplicationStateStarting:
		return "STARTING"
	case ApplicationStateRunning:
		return "RUNNING"
	case ApplicationStateErrored:
		return "ERRORED"
	default:
		return ""
	}
}

func ApplicationStateFromString(str string) ApplicationState {
	switch str {
	case "IDLE":
		return ApplicationStateIdle
	case "STARTING":
		return ApplicationStateStarting
	case "RUNNING":
		return ApplicationStateRunning
	case "ERRORED":
		return ApplicationStateErrored
	default:
		panic(fmt.Sprintf("unknown application state: %v", str))
	}
}

type Application struct {
	ID            string
	Repository    Repository
	BranchName    string
	BuildType     builder.BuildType
	State         ApplicationState
	CurrentCommit string
	WantCommit    string
}

type Artifact struct {
	ID        string
	Size      int64
	CreatedAt time.Time
}

type Build struct {
	ID            string
	Commit        string
	Status        builder.BuildStatus
	ApplicationID string
	StartedAt     time.Time
	FinishedAt    optional.Of[time.Time]
	Artifact      optional.Of[Artifact]
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
