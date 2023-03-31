package domain

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

import (
	"context"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type GetApplicationCondition struct {
	IDIn      optional.Of[[]string]
	UserID    optional.Of[string]
	BuildType optional.Of[builder.BuildType]
	State     optional.Of[ApplicationState]
	// InSync WantCommit が CurrentCommit に一致する
	InSync optional.Of[bool]
}

type UpdateApplicationArgs struct {
	State         optional.Of[ApplicationState]
	CurrentCommit optional.Of[string]
	WantCommit    optional.Of[string]
}

type ApplicationRepository interface {
	GetApplications(ctx context.Context, cond GetApplicationCondition) ([]*Application, error)
	GetApplication(ctx context.Context, id string) (*Application, error)
	CreateApplication(ctx context.Context, app *Application) error
	UpdateApplication(ctx context.Context, id string, args UpdateApplicationArgs) error
	RegisterApplicationOwner(ctx context.Context, applicationID string, userID string) error
	GetWebsites(ctx context.Context, applicationIDs []string) ([]*Website, error)
	AddWebsite(ctx context.Context, applicationID string, website *Website) error
	DeleteWebsite(ctx context.Context, applicationID string, websiteID string) error
}

type ArtifactRepository interface {
	CreateArtifact(ctx context.Context, size int64, buildID string, sid string) error
}

type AvailableDomainRepository interface {
	GetAvailableDomains(ctx context.Context) (AvailableDomainSlice, error)
	AddAvailableDomain(ctx context.Context, ad *AvailableDomain) error
	DeleteAvailableDomain(ctx context.Context, domain string) error
}

type GetBuildCondition struct {
	ApplicationID optional.Of[string]
	Commit        optional.Of[string]
	CommitIn      optional.Of[[]string]
	Status        optional.Of[builder.BuildStatus]
	Retriable     optional.Of[bool]
}

type UpdateBuildArgs struct {
	FromStatus optional.Of[builder.BuildStatus]
	Status     optional.Of[builder.BuildStatus]
	StartedAt  optional.Of[time.Time]
	UpdatedAt  optional.Of[time.Time]
	FinishedAt optional.Of[time.Time]
}

type BuildRepository interface {
	GetBuilds(ctx context.Context, condition GetBuildCondition) ([]*Build, error)
	GetBuild(ctx context.Context, buildID string) (*Build, error)
	CreateBuild(ctx context.Context, build *Build) error
	UpdateBuild(ctx context.Context, id string, args UpdateBuildArgs) error
	MarkCommitAsRetriable(ctx context.Context, applicationID string, commit string) error
}

type EnvironmentRepository interface {
	GetEnv(ctx context.Context, applicationID string) ([]*Environment, error)
	SetEnv(ctx context.Context, applicationID, key, value string) error
}

type GetRepositoryCondition struct {
	UserID optional.Of[string]
}

type GitRepositoryRepository interface {
	GetRepositories(ctx context.Context, condition GetRepositoryCondition) ([]*Repository, error)
	GetRepository(ctx context.Context, id string) (*Repository, error)
	CreateRepository(ctx context.Context, repo *Repository) error
	RegisterRepositoryOwner(ctx context.Context, repositoryID string, userID string) error
}

type CreateUserArgs struct {
	Name string
}

type UserRepository interface {
	CreateUser(ctx context.Context, args CreateUserArgs) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}
