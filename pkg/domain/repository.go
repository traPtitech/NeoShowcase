package domain

import (
	"context"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type GetApplicationCondition struct {
	IDIn         optional.Of[[]string]
	RepositoryID optional.Of[string]
	UserID       optional.Of[string]
	DeployType   optional.Of[DeployType]
	Running      optional.Of[bool]
}

type UpdateApplicationArgs struct {
	Name             optional.Of[string]
	RepositoryID     optional.Of[string]
	RefName          optional.Of[string]
	Commit           optional.Of[string]
	Running          optional.Of[bool]
	Container        optional.Of[ContainerState]
	ContainerMessage optional.Of[string]
	CurrentBuild     optional.Of[string]
	UpdatedAt        optional.Of[time.Time]
	Config           optional.Of[ApplicationConfig]
	Websites         optional.Of[[]*Website]
	PortPublications optional.Of[[]*PortPublication]
	OwnerIDs         optional.Of[[]string]
}

func (a *Application) Apply(args *UpdateApplicationArgs) {
	if args.Name.Valid {
		a.Name = args.Name.V
	}
	if args.RepositoryID.Valid {
		a.RepositoryID = args.RepositoryID.V
	}
	if args.RefName.Valid {
		a.RefName = args.RefName.V
	}
	if args.Commit.Valid {
		a.Commit = args.Commit.V
	}
	if args.Running.Valid {
		a.Running = args.Running.V
	}
	if args.Container.Valid {
		a.Container = args.Container.V
	}
	if args.ContainerMessage.Valid {
		a.ContainerMessage = args.ContainerMessage.V
	}
	if args.CurrentBuild.Valid {
		a.CurrentBuild = args.CurrentBuild.V
	}
	if args.UpdatedAt.Valid {
		a.UpdatedAt = args.UpdatedAt.V
	}
	if args.Config.Valid {
		a.DeployType = args.Config.V.BuildConfig.BuildType().DeployType()
		a.Config = args.Config.V
	}
	if args.Websites.Valid {
		a.Websites = args.Websites.V
	}
	if args.PortPublications.Valid {
		a.PortPublications = args.PortPublications.V
	}
	if args.OwnerIDs.Valid {
		a.OwnerIDs = args.OwnerIDs.V
	}
}

type ApplicationRepository interface {
	GetApplications(ctx context.Context, cond GetApplicationCondition) ([]*Application, error)
	GetApplication(ctx context.Context, id string) (*Application, error)
	CreateApplication(ctx context.Context, app *Application) error
	UpdateApplication(ctx context.Context, id string, args *UpdateApplicationArgs) error
	BulkUpdateState(ctx context.Context, states []*Container) error
	DeleteApplication(ctx context.Context, id string) error
}

type GetArtifactCondition struct {
	ApplicationID optional.Of[string]
	IsDeleted     optional.Of[bool]
}

type UpdateArtifactArgs struct {
	DeletedAt optional.Of[time.Time]
}

type ArtifactRepository interface {
	GetArtifact(ctx context.Context, id string) (*Artifact, error)
	GetArtifacts(ctx context.Context, cond GetArtifactCondition) ([]*Artifact, error)
	CreateArtifact(ctx context.Context, artifact *Artifact) error
	UpdateArtifact(ctx context.Context, id string, args UpdateArtifactArgs) error
	HardDeleteArtifacts(ctx context.Context, cond GetArtifactCondition) error
}

type GetBuildCondition struct {
	ID            optional.Of[string]
	IDIn          optional.Of[[]string]
	ApplicationID optional.Of[string]
	Commit        optional.Of[string]
	CommitIn      optional.Of[[]string]
	ConfigHash    optional.Of[string]
	Status        optional.Of[BuildStatus]
	Retriable     optional.Of[bool]
	Offset        optional.Of[int]
	Limit         optional.Of[int]
	SortAsc       optional.Of[bool]
}

type UpdateBuildArgs struct {
	Status     optional.Of[BuildStatus]
	StartedAt  optional.Of[time.Time]
	UpdatedAt  optional.Of[time.Time]
	FinishedAt optional.Of[time.Time]
}

type BuildRepository interface {
	GetBuilds(ctx context.Context, cond GetBuildCondition) ([]*Build, error)
	GetLatestBuilds(ctx context.Context, appIDIn []string) ([]*Build, error)
	GetBuild(ctx context.Context, buildID string) (*Build, error)
	CreateBuild(ctx context.Context, build *Build) error
	UpdateBuild(ctx context.Context, cond GetBuildCondition, args UpdateBuildArgs) (int64, error)
	MarkCommitAsRetriable(ctx context.Context, applicationID string, commit string) error
	DeleteBuilds(ctx context.Context, cond GetBuildCondition) error
}

type GetEnvCondition struct {
	ApplicationIDIn optional.Of[[]string]
	ApplicationID   optional.Of[string]
	Key             optional.Of[string]
}

type EnvironmentRepository interface {
	GetEnv(ctx context.Context, cond GetEnvCondition) ([]*Environment, error)
	SetEnv(ctx context.Context, env *Environment) error
	DeleteEnv(ctx context.Context, cond GetEnvCondition) error
}

type GetRepositoryCondition struct {
	IDs                optional.Of[[]string]
	URLs               optional.Of[[]string]
	CreatableOrOwnedBy optional.Of[string]
	PublicOrOwnedBy    optional.Of[string]
	UserID             optional.Of[string]
}

type UpdateRepositoryArgs struct {
	Name     optional.Of[string]
	URL      optional.Of[string]
	Auth     optional.Of[optional.Of[RepositoryAuth]]
	OwnerIDs optional.Of[[]string]
}

func (r *Repository) Apply(args *UpdateRepositoryArgs) {
	if args.Name.Valid {
		r.Name = args.Name.V
	}
	if args.URL.Valid {
		r.URL = args.URL.V
	}
	if args.Auth.Valid {
		r.Auth = args.Auth.V
	}
	if args.OwnerIDs.Valid {
		r.OwnerIDs = args.OwnerIDs.V
	}
}

type GitRepositoryRepository interface {
	GetRepositories(ctx context.Context, condition GetRepositoryCondition) ([]*Repository, error)
	GetRepository(ctx context.Context, id string) (*Repository, error)
	CreateRepository(ctx context.Context, repo *Repository) error
	UpdateRepository(ctx context.Context, id string, args *UpdateRepositoryArgs) error
	DeleteRepository(ctx context.Context, id string) error
}

type RepositoryCommitRepository interface {
	GetCommits(ctx context.Context, hashes []string) ([]*RepositoryCommit, error)
	RecordCommit(ctx context.Context, commit *RepositoryCommit) error
}

type CreateUserArgs struct {
	Name string
}

type GetUserCondition struct {
	Admin optional.Of[bool]
}

type GetUserKeyCondition struct {
	UserIDs optional.Of[[]string]
}

type UserRepository interface {
	EnsureUser(ctx context.Context, name string) (*User, error)
	EnsureUsers(ctx context.Context, names []string) ([]*User, error)
	GetUsers(ctx context.Context, cond GetUserCondition) ([]*User, error)
	GetUserKeys(ctx context.Context, cond GetUserKeyCondition) ([]*UserKey, error)
	CreateUserKey(ctx context.Context, key *UserKey) error
	DeleteUserKey(ctx context.Context, keyID string, userID string) error
}
