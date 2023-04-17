package usecase

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

func handleRepoError[T any](entity T, err error) (T, error) {
	switch err {
	case repository.ErrNotFound:
		return entity, newError(ErrorTypeNotFound, "not found", err)
	default:
		return entity, err
	}
}

type APIServerService struct {
	bus             domain.Bus
	artifactRepo    domain.ArtifactRepository
	appRepo         domain.ApplicationRepository
	adRepo          domain.AvailableDomainRepository
	buildRepo       domain.BuildRepository
	envRepo         domain.EnvironmentRepository
	gitRepo         domain.GitRepositoryRepository
	storage         domain.Storage
	mariaDBManager  domain.MariaDBManager
	mongoDBManager  domain.MongoDBManager
	containerLogger domain.ContainerLogger
	controller      domain.ControllerServiceClient
}

func NewAPIServerService(
	bus domain.Bus,
	artifactRepo domain.ArtifactRepository,
	appRepo domain.ApplicationRepository,
	adRepo domain.AvailableDomainRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
	storage domain.Storage,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
	containerLogger domain.ContainerLogger,
	controller domain.ControllerServiceClient,
) *APIServerService {
	return &APIServerService{
		bus:             bus,
		artifactRepo:    artifactRepo,
		appRepo:         appRepo,
		adRepo:          adRepo,
		buildRepo:       buildRepo,
		envRepo:         envRepo,
		gitRepo:         gitRepo,
		storage:         storage,
		mariaDBManager:  mariaDBManager,
		mongoDBManager:  mongoDBManager,
		containerLogger: containerLogger,
		controller:      controller,
	}
}

func (s *APIServerService) isRepositoryOwner(ctx context.Context, id string) error {
	user := web.GetUser(ctx)
	if user.Admin {
		return nil
	}
	repo, err := s.gitRepo.GetRepository(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get repository")
	}
	if !lo.Contains(repo.OwnerIDs, user.ID) {
		return newError(ErrorTypeForbidden, "you do not have permission for this repository", nil)
	}
	return nil
}

func (s *APIServerService) isApplicationOwner(ctx context.Context, id string) error {
	user := web.GetUser(ctx)
	if user.Admin {
		return nil
	}
	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get application")
	}
	if !lo.Contains(app.OwnerIDs, user.ID) {
		return newError(ErrorTypeForbidden, "you do not have permission for this application", nil)
	}
	return nil
}

func (s *APIServerService) isBuildOwner(ctx context.Context, id string) error {
	user := web.GetUser(ctx)
	if user.Admin {
		return nil
	}
	build, err := s.buildRepo.GetBuild(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get build")
	}
	app, err := s.appRepo.GetApplication(ctx, build.ApplicationID)
	if err != nil {
		return errors.Wrap(err, "failed to get application")
	}
	if !lo.Contains(app.OwnerIDs, user.ID) {
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
