package apiserver

import (
	"context"
	"github.com/motoki317/sc"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/util/scutil"
)

func handleRepoError[T any](entity T, err error) (T, error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return entity, newError(ErrorTypeNotFound, "not found", err)
	default:
		return entity, err
	}
}

type Service struct {
	artifactRepo    domain.ArtifactRepository
	appRepo         domain.ApplicationRepository
	buildRepo       domain.BuildRepository
	envRepo         domain.EnvironmentRepository
	gitRepo         domain.GitRepositoryRepository
	commitRepo      domain.RepositoryCommitRepository
	userRepo        domain.UserRepository
	storage         domain.Storage
	mariaDBManager  domain.MariaDBManager
	mongoDBManager  domain.MongoDBManager
	metricsService  domain.MetricsService
	containerLogger domain.ContainerLogger
	controller      domain.ControllerServiceClient
	fallbackKey     *ssh.PublicKeys
	image           builder.ImageConfig

	systemInfo *sc.Cache[struct{}, *domain.SystemInfo]
	tmpKeys    *tmpKeyPairService
}

func NewService(
	artifactRepo domain.ArtifactRepository,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
	commitRepo domain.RepositoryCommitRepository,
	userRepo domain.UserRepository,
	storage domain.Storage,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
	metricsService domain.MetricsService,
	containerLogger domain.ContainerLogger,
	controller domain.ControllerServiceClient,
	image builder.ImageConfig,
	fallbackKey *ssh.PublicKeys,
) (*Service, error) {
	return &Service{
		artifactRepo:    artifactRepo,
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		envRepo:         envRepo,
		gitRepo:         gitRepo,
		commitRepo:      commitRepo,
		userRepo:        userRepo,
		storage:         storage,
		mariaDBManager:  mariaDBManager,
		mongoDBManager:  mongoDBManager,
		metricsService:  metricsService,
		containerLogger: containerLogger,
		controller:      controller,
		fallbackKey:     fallbackKey,
		image:           image,

		systemInfo: sc.NewMust(scutil.WrapFunc(controller.GetSystemInfo), 5*time.Minute, 10*time.Minute),
		tmpKeys:    newTmpKeyPairService(),
	}, nil
}

func (s *Service) isRepositoryOwner(ctx context.Context, repoID string) error {
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

func (s *Service) isApplicationOwner(ctx context.Context, appID string) error {
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

func (s *Service) isBuildOwner(ctx context.Context, buildID string) error {
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

func (s *Service) isAdmin(ctx context.Context) error {
	user := web.GetUser(ctx)
	if !user.Admin {
		return newError(ErrorTypeForbidden, "you do not have permission for this action", nil)
	}
	return nil
}
