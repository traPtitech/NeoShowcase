package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/random"
)

func handleRepoError[T any](entity T, err error) (T, error) {
	switch err {
	case repository.ErrNotFound:
		return entity, newError(ErrorTypeNotFound, "not found", err)
	default:
		return entity, err
	}
}

type CreateApplicationArgs struct {
	UserID        string
	Name          string
	RepositoryURL string
	BranchName    string
	BuildType     builder.BuildType
	Config        domain.ApplicationConfig
	Websites      []*domain.Website
	StartOnCreate bool
}

type APIServerService interface {
	GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error)
	GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error)
	AddAvailableDomain(ctx context.Context, domain string) error
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	GetApplication(ctx context.Context, id string) (*domain.Application, error)
	DeleteApplication(ctx context.Context, id string) error
	GetApplicationBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error)
	GetApplicationBuild(ctx context.Context, buildID string) (*domain.Build, error)
	SetApplicationEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error
	GetApplicationEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error)
	CancelBuild(ctx context.Context, buildID string) error
	RetryCommitBuild(ctx context.Context, applicationID string, commit string) error
	StartApplication(ctx context.Context, id string) error
	StopApplication(ctx context.Context, id string) error
}

type apiServerService struct {
	bus            domain.Bus
	appRepo        domain.ApplicationRepository
	adRepo         domain.AvailableDomainRepository
	buildRepo      domain.BuildRepository
	envRepo        domain.EnvironmentRepository
	gitRepo        domain.GitRepositoryRepository
	deploySvc      AppDeployService
	backend        domain.Backend
	component      domain.ComponentService
	mariaDBManager domain.MariaDBManager
	mongoDBManager domain.MongoDBManager
}

func NewAPIServerService(
	bus domain.Bus,
	appRepo domain.ApplicationRepository,
	adRepo domain.AvailableDomainRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
	deploySvc AppDeployService,
	backend domain.Backend,
	component domain.ComponentService,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
) APIServerService {
	return &apiServerService{
		bus:            bus,
		appRepo:        appRepo,
		adRepo:         adRepo,
		buildRepo:      buildRepo,
		envRepo:        envRepo,
		gitRepo:        gitRepo,
		deploySvc:      deploySvc,
		backend:        backend,
		component:      component,
		mariaDBManager: mariaDBManager,
		mongoDBManager: mongoDBManager,
	}
}

func (s *apiServerService) GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error) {
	return s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{UserID: optional.From(userID)})
}

func (s *apiServerService) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	return s.adRepo.GetAvailableDomains(ctx)
}

func (s *apiServerService) AddAvailableDomain(ctx context.Context, d string) error {
	ad := domain.AvailableDomain{Domain: d}
	if !ad.IsValid() {
		return newError(ErrorTypeBadRequest, "invalid new domain", nil)
	}
	return s.adRepo.AddAvailableDomain(ctx, d)
}

func (s *apiServerService) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	repo, err := s.gitRepo.GetRepository(ctx, args.RepositoryURL)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	if err == repository.ErrNotFound {
		repoName, err := domain.ExtractNameFromRepositoryURL(args.RepositoryURL)
		if err != nil {
			return nil, newError(ErrorTypeBadRequest, "malformed repository url", err)
		}
		repo, err = s.gitRepo.RegisterRepository(ctx, domain.RegisterRepositoryArgs{
			Name: repoName,
			URL:  args.RepositoryURL,
		})
		if err != nil {
			return nil, err
		}
	}

	domains, err := s.adRepo.GetAvailableDomains(ctx)
	if err != nil {
		return nil, err
	}
	for _, website := range args.Websites {
		if !website.IsValid() {
			return nil, newError(ErrorTypeBadRequest, "invalid website", nil)
		}
		if !domains.Match(website.FQDN) {
			return nil, newError(ErrorTypeBadRequest, "domain not available", nil)
		}
	}

	var initialState domain.ApplicationState
	if args.StartOnCreate {
		initialState = domain.ApplicationStateDeploying
	} else {
		initialState = domain.ApplicationStateIdle
	}
	application, err := s.appRepo.CreateApplication(ctx, domain.CreateApplicationArgs{
		Name:         args.Name,
		RepositoryID: repo.ID,
		BranchName:   args.BranchName,
		BuildType:    args.BuildType,
		State:        initialState,
		Config:       args.Config,
		Websites:     args.Websites,
	})
	if err != nil {
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

	s.bus.Publish(event.FetcherFetchRequest, nil)

	return application, nil
}

func (s *apiServerService) createApplicationDatabase(ctx context.Context, app *domain.Application) error {
	dbName := fmt.Sprintf("nsapp_%s", app.ID)

	if app.Config.UseMariaDB {
		dbPassword := random.SecureGeneratePassword(32)
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

	if app.Config.UseMongoDB {
		dbPassword := random.SecureGeneratePassword(32)
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
	application, err := s.appRepo.GetApplication(ctx, id)
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
	return s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(applicationID)})
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

func (s *apiServerService) CancelBuild(_ context.Context, buildID string) error {
	s.component.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_CANCEL_BUILD,
		Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.CancelBuildRequest{BuildId: buildID}},
	})
	return nil
}

func (s *apiServerService) RetryCommitBuild(ctx context.Context, applicationID string, commit string) error {
	err := s.buildRepo.MarkCommitAsRetriable(ctx, applicationID, commit)
	if err != nil {
		return err
	}
	s.bus.Publish(event.CDServiceRegisterBuildRequest, nil)
	return nil
}

func (s *apiServerService) StartApplication(_ context.Context, id string) error {
	ok := s.deploySvc.Synchronize(id, true)
	if !ok {
		return errors.New("application is currently busy")
	}
	return nil
}

func (s *apiServerService) StopApplication(_ context.Context, id string) error {
	ok := s.deploySvc.Stop(id)
	if !ok {
		return errors.New("application is currently busy")
	}
	return nil
}
