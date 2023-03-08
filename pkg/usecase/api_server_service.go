package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

func handleRepoError[T any](entity T, err error) (T, error) {
	switch err {
	case repository.ErrNotFound:
		return entity, ErrNotFound
	default:
		return entity, err
	}
}

type CreateApplicationArgs struct {
	UserID        string
	RepositoryURL string
	BranchName    string
	BuildType     builder.BuildType
}

type APIServerService interface {
	GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error)
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	GetApplication(ctx context.Context, id string) (*domain.Application, error)
	DeleteApplication(ctx context.Context, id string) error
	GetApplicationBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error)
	GetApplicationBuild(ctx context.Context, buildID string) (*domain.Build, error)
	SetApplicationEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error
	GetApplicationEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error)
	StartApplication(ctx context.Context, id string) error
	RestartApplication(ctx context.Context, id string) error
	StopApplication(ctx context.Context, id string) error
}

type apiServerService struct {
	appRepo        repository.ApplicationRepository
	buildRepo      repository.BuildRepository
	envRepo        repository.EnvironmentRepository
	gitRepo        repository.GitRepositoryRepository
	deploySvc      AppDeployService
	backend        domain.Backend
	mariaDBManager domain.MariaDBManager
	mongoDBManager domain.MongoDBManager
}

func NewAPIServerService(
	appRepo repository.ApplicationRepository,
	buildRepo repository.BuildRepository,
	envRepo repository.EnvironmentRepository,
	gitRepo repository.GitRepositoryRepository,
	deploySvc AppDeployService,
	backend domain.Backend,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
) APIServerService {
	return &apiServerService{
		appRepo:        appRepo,
		buildRepo:      buildRepo,
		envRepo:        envRepo,
		gitRepo:        gitRepo,
		deploySvc:      deploySvc,
		backend:        backend,
		mariaDBManager: mariaDBManager,
		mongoDBManager: mongoDBManager,
	}
}

func (s *apiServerService) GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error) {
	return s.appRepo.GetApplicationsByUserID(ctx, userID)
}

func (s *apiServerService) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	repo, err := s.gitRepo.GetRepository(ctx, args.RepositoryURL)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	if err == repository.ErrNotFound {
		repoName, err := domain.ExtractNameFromRepositoryURL(args.RepositoryURL)
		if err != nil {
			return nil, fmt.Errorf("malformed repository url: %w", err)
		}
		repo, err = s.gitRepo.RegisterRepository(ctx, repository.RegisterRepositoryArgs{
			Name: repoName,
			URL:  args.RepositoryURL,
		})
		if err != nil {
			return nil, err
		}
	}

	application, err := s.appRepo.CreateApplication(ctx, repository.CreateApplicationArgs{
		RepositoryID: repo.ID,
		BranchName:   args.BranchName,
		BuildType:    args.BuildType,
	})
	if err != nil {
		if err == repository.ErrDuplicate {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}

	err = s.appRepo.RegisterApplicationOwner(ctx, application.ID, args.UserID)
	if err != nil {
		return nil, err
	}

	err = s.createApplicationDatabase(ctx, application)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (s *apiServerService) createApplicationDatabase(ctx context.Context, app *domain.Application) error {
	dbName := fmt.Sprintf("%s_%s", app.Repository, app.ID)

	// TODO: アプリケーションの設定の取得
	applicationNeedsMariaDB := true
	if applicationNeedsMariaDB {
		dbPassword := generateRandomString(32)
		dbSetting := domain.CreateArgs{
			Database: dbName,
			Password: dbPassword,
		}
		if err := s.mariaDBManager.Create(ctx, dbSetting); err != nil {
			return err
		}

		if err := s.envRepo.SetEnv(ctx, app.ID, domain.EnvMySQLUserKey, dbName); err != nil {
			return err
		}
		if err := s.envRepo.SetEnv(ctx, app.ID, domain.EnvMySQLPasswordKey, dbPassword); err != nil {
			return err
		}
		if err := s.envRepo.SetEnv(ctx, app.ID, domain.EnvMySQLDatabaseKey, dbName); err != nil {
			return err
		}
	}

	// TODO: アプリケーションの設定の取得
	applicationNeedsMongoDB := true
	if applicationNeedsMongoDB {
		dbPassword := generateRandomString(32)
		dbSetting := domain.CreateArgs{
			Database: dbName,
			Password: dbPassword,
		}
		err := s.mongoDBManager.Create(ctx, dbSetting)
		if err != nil {
			return err
		}

		if err := s.envRepo.SetEnv(ctx, app.ID, domain.EnvMongoDBUserKey, dbName); err != nil {
			return err
		}
		if err := s.envRepo.SetEnv(ctx, app.ID, domain.EnvMongoDBPasswordKey, dbPassword); err != nil {
			return err
		}
		if err := s.envRepo.SetEnv(ctx, app.ID, domain.EnvMongoDBDatabaseKey, dbName); err != nil {
			return err
		}
	}

	return nil
}

func (s *apiServerService) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	application, err := s.appRepo.GetApplicationByID(ctx, id)
	return handleRepoError(application, err)
}

func (s *apiServerService) DeleteApplication(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
	// delete artifacts
	// delete builds
	// delete websites
	// delete environments
	// delete owners
	// s.deleteApplicationDatabase()
}

func (s *apiServerService) deleteApplicationDatabase(ctx context.Context, app *domain.Application) error {
	// TODO implement me
	panic("implement me")
}

func (s *apiServerService) GetApplicationBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, applicationID)
}

func (s *apiServerService) GetApplicationBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	return handleRepoError(build, err)
}

func (s *apiServerService) GetApplicationEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	return s.envRepo.GetEnv(ctx, applicationID)
}

func (s *apiServerService) SetApplicationEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error {
	return s.envRepo.SetEnv(ctx, applicationID, key, value)
}

func (s *apiServerService) StartApplication(ctx context.Context, id string) error {
	build, err := s.buildRepo.GetLastSuccessBuild(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return s.deploySvc.QueueDeployment(ctx, id, build.ID)
}

func (s *apiServerService) RestartApplication(ctx context.Context, id string) error {
	return s.backend.RestartContainer(ctx, id)
}

func (s *apiServerService) StopApplication(ctx context.Context, id string) error {
	return s.backend.DestroyContainer(ctx, id)
}
