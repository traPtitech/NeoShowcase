package domain

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type GetApplicationCondition struct {
	UserID    optional.Of[string]
	BuildType optional.Of[builder.BuildType]
	State     optional.Of[ApplicationState]
	// InSync WantCommit が CurrentCommit に一致する
	InSync optional.Of[bool]
}

type CreateWebsiteArgs struct {
	FQDN       string
	PathPrefix string
	HTTPS      bool
	HTTPPort   int
}

type CreateApplicationArgs struct {
	Name         string
	RepositoryID string
	BranchName   string
	BuildType    builder.BuildType
	State        ApplicationState
	Config       ApplicationConfig
	Websites     []*CreateWebsiteArgs
}

type UpdateApplicationArgs struct {
	State         optional.Of[ApplicationState]
	CurrentCommit optional.Of[string]
	WantCommit    optional.Of[string]
}

type ApplicationRepository interface {
	GetApplications(ctx context.Context, cond GetApplicationCondition) ([]*Application, error)
	GetApplication(ctx context.Context, id string) (*Application, error)
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*Application, error)
	UpdateApplication(ctx context.Context, id string, args UpdateApplicationArgs) error
	RegisterApplicationOwner(ctx context.Context, applicationID string, userID string) error
	GetWebsites(ctx context.Context, applicationIDs []string) ([]*Website, error)
	AddWebsite(ctx context.Context, applicationID string, args CreateWebsiteArgs) error
	DeleteWebsite(ctx context.Context, applicationID string, websiteID string) error
}

type ArtifactRepository interface {
	CreateArtifact(ctx context.Context, size int64, buildID string, sid string) error
}

type AvailableDomainRepository interface {
	GetAvailableDomains(ctx context.Context) (AvailableDomainSlice, error)
	AddAvailableDomain(ctx context.Context, domain string) error
	DeleteAvailableDomain(ctx context.Context, domain string) error
}

type UpdateBuildArgs struct {
	ID     string
	Status builder.BuildStatus
}

type BuildRepository interface {
	GetBuilds(ctx context.Context, applicationID string) ([]*Build, error)
	GetBuildsInCommit(ctx context.Context, commits []string) ([]*Build, error)
	GetBuild(ctx context.Context, buildID string) (*Build, error)
	CreateBuild(ctx context.Context, applicationID string, commit string) (*Build, error)
	UpdateBuild(ctx context.Context, args UpdateBuildArgs) error
	MarkCommitAsRetriable(ctx context.Context, applicationID string, commit string) error
}

type EnvironmentRepository interface {
	GetEnv(ctx context.Context, applicationID string) ([]*Environment, error)
	SetEnv(ctx context.Context, applicationID, key, value string) error
}

type RegisterRepositoryArgs struct {
	Name string
	URL  string
}

type GitRepositoryRepository interface {
	RegisterRepository(ctx context.Context, args RegisterRepositoryArgs) (Repository, error)
	GetRepositoryByID(ctx context.Context, id string) (Repository, error)
	GetRepository(ctx context.Context, rawURL string) (Repository, error)
}

type CreateUserArgs struct {
	Name string
}

type UserRepository interface {
	CreateUser(ctx context.Context, args CreateUserArgs) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}
