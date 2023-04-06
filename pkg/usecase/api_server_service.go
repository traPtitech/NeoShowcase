package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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

type APIServerService interface {
	GetRepositories(ctx context.Context) ([]*domain.Repository, error)
	CreateRepository(ctx context.Context, repo *domain.Repository) error
	GetApplications(ctx context.Context) ([]*domain.Application, error)
	GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error)
	AddAvailableDomain(ctx context.Context, ad *domain.AvailableDomain) error
	CreateApplication(ctx context.Context, app *domain.Application) (*domain.Application, error)
	GetApplication(ctx context.Context, id string) (*domain.Application, error)
	UpdateApplication(ctx context.Context, app *domain.Application, args *domain.UpdateApplicationArgs) error
	DeleteApplication(ctx context.Context, id string) error
	GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error)
	GetBuild(ctx context.Context, buildID string) (*domain.Build, error)
	GetBuildLog(ctx context.Context, buildID string) ([]byte, error)
	GetArtifact(ctx context.Context, artifactID string) ([]byte, error)
	SetEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error
	GetEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error)
	CancelBuild(ctx context.Context, buildID string) error
	RetryCommitBuild(ctx context.Context, applicationID string, commit string) error
	StartApplication(ctx context.Context, id string) error
	StopApplication(ctx context.Context, id string) error
}

type apiServerService struct {
	bus            domain.Bus
	artifactRepo   domain.ArtifactRepository
	appRepo        domain.ApplicationRepository
	adRepo         domain.AvailableDomainRepository
	buildRepo      domain.BuildRepository
	envRepo        domain.EnvironmentRepository
	gitRepo        domain.GitRepositoryRepository
	deploySvc      AppDeployService
	backend        domain.Backend
	storage        domain.Storage
	component      domain.ComponentService
	mariaDBManager domain.MariaDBManager
	mongoDBManager domain.MongoDBManager
}

func NewAPIServerService(
	bus domain.Bus,
	artifactRepo domain.ArtifactRepository,
	appRepo domain.ApplicationRepository,
	adRepo domain.AvailableDomainRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
	deploySvc AppDeployService,
	backend domain.Backend,
	storage domain.Storage,
	component domain.ComponentService,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
) APIServerService {
	return &apiServerService{
		bus:            bus,
		artifactRepo:   artifactRepo,
		appRepo:        appRepo,
		adRepo:         adRepo,
		buildRepo:      buildRepo,
		envRepo:        envRepo,
		gitRepo:        gitRepo,
		deploySvc:      deploySvc,
		backend:        backend,
		storage:        storage,
		component:      component,
		mariaDBManager: mariaDBManager,
		mongoDBManager: mongoDBManager,
	}
}

func (s *apiServerService) GetRepositories(ctx context.Context) ([]*domain.Repository, error) {
	return s.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{})
}

func (s *apiServerService) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	return s.gitRepo.CreateRepository(ctx, repo)
}

func (s *apiServerService) GetApplications(ctx context.Context) ([]*domain.Application, error) {
	return s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
}

func (s *apiServerService) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	return s.adRepo.GetAvailableDomains(ctx)
}

func (s *apiServerService) AddAvailableDomain(ctx context.Context, ad *domain.AvailableDomain) error {
	if !ad.IsValid() {
		return newError(ErrorTypeBadRequest, "invalid new domain", nil)
	}
	return s.adRepo.AddAvailableDomain(ctx, ad)
}

func (s *apiServerService) CreateApplication(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	domains, err := s.adRepo.GetAvailableDomains(ctx)
	if err != nil {
		return nil, err
	}
	for _, website := range app.Websites {
		if !website.IsValid() {
			return nil, newError(ErrorTypeBadRequest, "invalid website", nil)
		}
		if !domains.IsAvailable(website.FQDN) {
			return nil, newError(ErrorTypeBadRequest, "domain not available", nil)
		}
	}

	err = s.appRepo.CreateApplication(ctx, app)
	if err != nil {
		return nil, err
	}

	err = s.createApplicationDatabase(ctx, app)
	if err != nil {
		return nil, err
	}

	s.bus.Publish(event.FetcherFetchRequest, nil)

	return s.GetApplication(ctx, app.ID)
}

func (s *apiServerService) createApplicationDatabase(ctx context.Context, app *domain.Application) error {
	dbName := domain.DBName(app.ID)

	if app.Config.UseMariaDB {
		host, port := s.mariaDBManager.GetHost()
		dbPassword := random.SecureGeneratePassword(32)
		dbSetting := domain.CreateArgs{
			Database: dbName,
			Password: dbPassword,
		}
		err := s.mariaDBManager.Create(ctx, dbSetting)
		if err != nil {
			return err
		}

		envs := []*domain.Environment{
			{ApplicationID: app.ID, Key: domain.EnvMySQLHostnameKey, Value: host, System: true},
			{ApplicationID: app.ID, Key: domain.EnvMySQLPortKey, Value: strconv.Itoa(port), System: true},
			{ApplicationID: app.ID, Key: domain.EnvMySQLUserKey, Value: dbName, System: true},
			{ApplicationID: app.ID, Key: domain.EnvMySQLPasswordKey, Value: dbPassword, System: true},
			{ApplicationID: app.ID, Key: domain.EnvMySQLDatabaseKey, Value: dbName, System: true},
		}
		for _, env := range envs {
			err = s.envRepo.SetEnv(ctx, env)
			if err != nil {
				return err
			}
		}
	}

	if app.Config.UseMongoDB {
		host, port := s.mongoDBManager.GetHost()
		dbPassword := random.SecureGeneratePassword(32)
		dbSetting := domain.CreateArgs{
			Database: dbName,
			Password: dbPassword,
		}
		err := s.mongoDBManager.Create(ctx, dbSetting)
		if err != nil {
			return err
		}

		envs := []*domain.Environment{
			{ApplicationID: app.ID, Key: domain.EnvMongoDBHostnameKey, Value: host, System: true},
			{ApplicationID: app.ID, Key: domain.EnvMongoDBPortKey, Value: strconv.Itoa(port), System: true},
			{ApplicationID: app.ID, Key: domain.EnvMongoDBUserKey, Value: dbName, System: true},
			{ApplicationID: app.ID, Key: domain.EnvMongoDBPasswordKey, Value: dbPassword, System: true},
			{ApplicationID: app.ID, Key: domain.EnvMongoDBDatabaseKey, Value: dbName, System: true},
		}
		for _, env := range envs {
			err = s.envRepo.SetEnv(ctx, env)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *apiServerService) deleteApplicationDatabase(ctx context.Context, app *domain.Application) error {
	dbName := domain.DBName(app.ID)

	if app.Config.UseMariaDB {
		err := s.mariaDBManager.Delete(ctx, domain.DeleteArgs{Database: dbName})
		if err != nil {
			return err
		}
	}

	if app.Config.UseMongoDB {
		err := s.mongoDBManager.Delete(ctx, domain.DeleteArgs{Database: dbName})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *apiServerService) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	application, err := s.appRepo.GetApplication(ctx, id)
	return handleRepoError(application, err)
}

func (s *apiServerService) UpdateApplication(ctx context.Context, app *domain.Application, args *domain.UpdateApplicationArgs) error {
	err := s.appRepo.UpdateApplication(ctx, app.ID, args)
	if err != nil {
		return err
	}
	return s.RetryCommitBuild(ctx, app.ID, app.CurrentCommit)
}

func (s *apiServerService) DeleteApplication(ctx context.Context, id string) error {
	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return err
	}
	if app.Running {
		return newError(ErrorTypeBadRequest, "stop the application first before deleting", nil)
	}

	err = s.deleteApplicationDatabase(ctx, app)
	if err != nil {
		return err
	}

	// delete artifacts
	artifacts, err := s.artifactRepo.GetArtifacts(ctx, domain.GetArtifactCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return err
	}
	for _, artifact := range artifacts {
		if artifact.DeletedAt.Valid {
			continue
		}
		err = domain.DeleteArtifact(s.storage, artifact.ID)
		if err != nil {
			log.Errorf("failed to delete artifact: %+v", err) // fail-safe
		}
	}
	err = s.artifactRepo.HardDeleteArtifacts(ctx, domain.GetArtifactCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return err
	}
	// delete builds
	builds, err := s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return err
	}
	for _, build := range builds {
		err = domain.DeleteBuildLog(s.storage, build.ID)
		if err != nil {
			log.Errorf("failed to delete build log: %+v", err) // fail-safe
		}
	}
	err = s.buildRepo.DeleteBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return err
	}
	// delete environments
	err = s.envRepo.DeleteEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return err
	}
	// delete websites, owners, application
	err = s.appRepo.DeleteApplication(ctx, app.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *apiServerService) GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(applicationID)})
}

func (s *apiServerService) GetBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	return handleRepoError(build, err)
}

func (s *apiServerService) GetBuildLog(ctx context.Context, buildID string) ([]byte, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		return nil, err
	}
	if !build.Status.IsFinished() {
		return nil, newError(ErrorTypeBadRequest, "build not finished", nil)
	}
	return domain.GetBuildLog(s.storage, buildID)
}

func (s *apiServerService) GetArtifact(_ context.Context, artifactID string) ([]byte, error) {
	return domain.GetArtifact(s.storage, artifactID)
}

func (s *apiServerService) GetEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	return s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(applicationID)})
}

func (s *apiServerService) SetEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error {
	env := &domain.Environment{ApplicationID: applicationID, Key: key, Value: value, System: false}
	return s.envRepo.SetEnv(ctx, env)
}

func (s *apiServerService) CancelBuild(_ context.Context, buildID string) error {
	s.component.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_CANCEL_BUILD,
		Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.BuildIdRequest{BuildId: buildID}},
	})
	return nil
}

func (s *apiServerService) RetryCommitBuild(ctx context.Context, applicationID string, commit string) error {
	err := s.buildRepo.MarkCommitAsRetriable(ctx, applicationID, commit)
	if err != nil {
		return err
	}
	// NOTE: requires the app to be running for builds to register
	s.bus.Publish(event.CDServiceRegisterBuildRequest, nil)
	return nil
}

func (s *apiServerService) StartApplication(ctx context.Context, id string) error {
	err := s.appRepo.UpdateApplication(ctx, id, &domain.UpdateApplicationArgs{
		Running:   optional.From(true),
		UpdatedAt: optional.From(time.Now()),
	})
	if err != nil {
		return errors.Wrap(err, "failed to mark application as running")
	}
	s.bus.Publish(event.CDServiceRegisterBuildRequest, nil)
	s.bus.Publish(event.CDServiceSyncDeployRequest, nil)
	return nil
}

func (s *apiServerService) StopApplication(ctx context.Context, id string) error {
	err := s.appRepo.UpdateApplication(ctx, id, &domain.UpdateApplicationArgs{
		Running:   optional.From(false),
		UpdatedAt: optional.From(time.Now()),
	})
	if err != nil {
		return errors.Wrap(err, "failed to mark application as not running")
	}
	s.bus.Publish(event.CDServiceSyncDeployRequest, nil)
	return nil
}
