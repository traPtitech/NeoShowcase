package apiserver

import (
	"context"
	"time"

	"github.com/motoki317/sc"

	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/scutil"
)

type Service struct {
	artifactRepo     domain.ArtifactRepository
	runtimeImageRepo domain.RuntimeImageRepository
	appRepo          domain.ApplicationRepository
	groupRepo        domain.ApplicationGroupRepository
	buildRepo        domain.BuildRepository
	envRepo          domain.EnvironmentRepository
	gitRepo          domain.GitRepositoryRepository
	commitRepo       domain.RepositoryCommitRepository
	userRepo         domain.UserRepository
	storage          domain.Storage
	mariaDBManager   domain.MariaDBManager
	mongoDBManager   domain.MongoDBManager
	metricsService   domain.MetricsService
	containerLogger  domain.ContainerLogger
	controller       domain.ControllerServiceClient
	gitsvc           domain.GitService
	regclient        builder.RegistryClient
	image            builder.ImageConfig

	systemInfo *sc.Cache[struct{}, *domain.SystemInfo]
	tmpKeys    *tmpKeyPairService
}

func NewService(
	artifactRepo domain.ArtifactRepository,
	runtimeImageRepo domain.RuntimeImageRepository,
	appRepo domain.ApplicationRepository,
	groupRepo domain.ApplicationGroupRepository,
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
	regclient builder.RegistryClient,
	image builder.ImageConfig,
	gitsvc domain.GitService,
) (*Service, error) {
	return &Service{
		artifactRepo:     artifactRepo,
		runtimeImageRepo: runtimeImageRepo,
		appRepo:          appRepo,
		groupRepo:        groupRepo,
		buildRepo:        buildRepo,
		envRepo:          envRepo,
		gitRepo:          gitRepo,
		commitRepo:       commitRepo,
		userRepo:         userRepo,
		storage:          storage,
		mariaDBManager:   mariaDBManager,
		mongoDBManager:   mongoDBManager,
		metricsService:   metricsService,
		containerLogger:  containerLogger,
		controller:       controller,
		gitsvc:           gitsvc,
		regclient:        regclient,
		image:            image,

		systemInfo: sc.NewMust(scutil.WrapFunc(controller.GetSystemInfo), 5*time.Minute, 10*time.Minute),
		tmpKeys:    newTmpKeyPairService(),
	}, nil
}

func (s *Service) isRepositoryOwner(ctx context.Context, repoID string) error {
	user := web.GetUser(ctx)
	repo, err := s.gitRepo.GetRepository(ctx, repoID)
	if err != nil {
		return oops.Wrapf(err, "getting repository")
	}
	if !repo.IsOwner(user) {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this repository", nil)
	}
	return nil
}

func (s *Service) isApplicationOwner(ctx context.Context, appID string) error {
	_, _, err := s.getOwnedApplication(ctx, appID)
	return err
}

// getOwnedApplication returns the application and its group when the user may operate the application.
func (s *Service) getOwnedApplication(ctx context.Context, appID string) (*domain.Application, *domain.ApplicationGroup, error) {
	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return nil, nil, oops.Wrapf(err, "getting application")
	}
	group, err := s.groupRepo.GetApplicationGroup(ctx, app.GroupID)
	if err != nil {
		return nil, nil, oops.Wrapf(err, "getting application group")
	}
	if err = checkGroupOwner(ctx, group); err != nil {
		return nil, nil, err
	}
	return app, group, nil
}

// checkGroupOwner checks that the user may operate the group and the applications in it.
func checkGroupOwner(ctx context.Context, group *domain.ApplicationGroup) error {
	if !group.IsOwner(web.GetUser(ctx)) {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this application group", nil)
	}
	return nil
}

func (s *Service) isBuildOwner(ctx context.Context, buildID string) error {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		return oops.Wrapf(err, "getting build")
	}
	return s.isApplicationOwner(ctx, build.ApplicationID)
}

func (s *Service) isAdmin(ctx context.Context) error {
	user := web.GetUser(ctx)
	if !user.Admin {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this action", nil)
	}
	return nil
}
