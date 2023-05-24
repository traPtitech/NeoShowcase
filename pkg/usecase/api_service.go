package usecase

import (
	"context"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
)

func handleRepoError[T any](entity T, err error) (T, error) {
	switch err {
	case repository.ErrNotFound:
		return entity, newError(ErrorTypeNotFound, "not found", err)
	default:
		return entity, err
	}
}

type AvatarBaseURL string

type APIServerService struct {
	avatarBaseURL   string
	artifactRepo    domain.ArtifactRepository
	appRepo         domain.ApplicationRepository
	buildRepo       domain.BuildRepository
	envRepo         domain.EnvironmentRepository
	gitRepo         domain.GitRepositoryRepository
	userRepo        domain.UserRepository
	storage         domain.Storage
	mariaDBManager  domain.MariaDBManager
	mongoDBManager  domain.MongoDBManager
	containerLogger domain.ContainerLogger
	controller      domain.ControllerServiceClient
}

func NewAPIServerService(
	avatarBaseURL AvatarBaseURL,
	artifactRepo domain.ArtifactRepository,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
	userRepo domain.UserRepository,
	storage domain.Storage,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
	containerLogger domain.ContainerLogger,
	controller domain.ControllerServiceClient,
) *APIServerService {
	return &APIServerService{
		avatarBaseURL:   string(avatarBaseURL),
		artifactRepo:    artifactRepo,
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		envRepo:         envRepo,
		gitRepo:         gitRepo,
		userRepo:        userRepo,
		storage:         storage,
		mariaDBManager:  mariaDBManager,
		mongoDBManager:  mongoDBManager,
		containerLogger: containerLogger,
		controller:      controller,
	}
}

func (s *APIServerService) isRepositoryOwner(ctx context.Context, repoID string) error {
	user := web.GetUser(ctx)
	repo, err := s.gitRepo.GetRepository(ctx, repoID)
	if err != nil {
		return errors.Wrap(err, "failed to get repository")
	}
	if !repo.IsOwner(user) {
		return newError(ErrorTypeForbidden, "you do not have permission for this repository", nil)
	}
	return nil
}

func (s *APIServerService) isApplicationOwner(ctx context.Context, appID string) error {
	user := web.GetUser(ctx)
	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return errors.Wrap(err, "failed to get application")
	}
	if !app.IsOwner(user) {
		return newError(ErrorTypeForbidden, "you do not have permission for this application", nil)
	}
	return nil
}

func (s *APIServerService) isBuildOwner(ctx context.Context, buildID string) error {
	user := web.GetUser(ctx)
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		return errors.Wrap(err, "failed to get build")
	}
	app, err := s.appRepo.GetApplication(ctx, build.ApplicationID)
	if err != nil {
		return errors.Wrap(err, "failed to get application")
	}
	if !app.IsOwner(user) {
		return newError(ErrorTypeForbidden, "you do not have permission for this application", nil)
	}
	return nil
}

func (s *APIServerService) isAdmin(ctx context.Context) error {
	user := web.GetUser(ctx)
	if !user.Admin {
		return newError(ErrorTypeForbidden, "you do not have permission for this action", nil)
	}
	return nil
}
