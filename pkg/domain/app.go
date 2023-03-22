package domain

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ApplicationState int

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateDeploying
	ApplicationStateRunning
	ApplicationStateErrored
)

func (s ApplicationState) String() string {
	switch s {
	case ApplicationStateIdle:
		return "IDLE"
	case ApplicationStateDeploying:
		return "DEPLOYING"
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
	case "DEPLOYING":
		return ApplicationStateDeploying
	case "RUNNING":
		return ApplicationStateRunning
	case "ERRORED":
		return ApplicationStateErrored
	default:
		panic(fmt.Sprintf("unknown application state: %v", str))
	}
}

type AuthenticationType int

const (
	AuthenticationTypeOff AuthenticationType = iota
	AuthenticationTypeSoft
	AuthenticationTypeHard
)

func (a AuthenticationType) String() string {
	switch a {
	case AuthenticationTypeOff:
		return "off"
	case AuthenticationTypeSoft:
		return "soft"
	case AuthenticationTypeHard:
		return "hard"
	default:
		return ""
	}
}

func AuthenticationTypeFromString(s string) AuthenticationType {
	switch s {
	case "off":
		return AuthenticationTypeOff
	case "soft":
		return AuthenticationTypeSoft
	case "hard":
		return AuthenticationTypeHard
	default:
		panic(fmt.Sprintf("unknown authentication type: %v", s))
	}
}

type ApplicationConfig struct {
	UseMariaDB     bool
	UseMongoDB     bool
	BaseImage      string
	DockerfileName string
	ArtifactPath   string
	BuildCmd       string
	EntrypointCmd  string
	Authentication AuthenticationType
}

var EmptyCommit = strings.Repeat("0", 40)

type Application struct {
	ID            string
	Name          string
	BranchName    string
	BuildType     builder.BuildType
	State         ApplicationState
	CurrentCommit string
	WantCommit    string

	Config     ApplicationConfig
	Repository Repository
	Websites   []*Website
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
	Retriable     bool
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

type Website struct {
	ID         string
	FQDN       string
	PathPrefix string
	HTTPS      bool
	HTTPPort   int
}

func GetActiveWebsites(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*Site, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		BuildType: optional.From(builder.BuildTypeStatic),
		State:     optional.From(ApplicationStateRunning),
	})
	if err != nil {
		return nil, err
	}

	commits := lo.Map(applications, func(app *Application, i int) string { return app.CurrentCommit })
	builds, err := buildRepo.GetBuildsInCommit(ctx, commits)
	if err != nil {
		return nil, err
	}

	// Last succeeded builds for each commit
	builds = lo.Filter(builds, func(build *Build, i int) bool { return build.Status == builder.BuildStatusSucceeded })
	slices.SortFunc(builds, func(a, b *Build) bool { return a.StartedAt.Before(b.StartedAt) })
	commitToBuild := lo.SliceToMap(builds, func(b *Build) (string, *Build) { return b.Commit, b })

	var sites []*Site
	for _, app := range applications {
		build, ok := commitToBuild[app.CurrentCommit]
		if !ok {
			continue
		}
		if !build.Artifact.Valid {
			continue
		}
		for _, website := range app.Websites {
			sites = append(sites, &Site{
				ID:            website.ID,
				FQDN:          website.FQDN,
				ArtifactID:    build.Artifact.V.ID,
				ApplicationID: app.ID,
			})
		}
	}
	return sites, nil
}
