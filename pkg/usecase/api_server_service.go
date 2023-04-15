package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
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

type APIServerService struct {
	bus             domain.Bus
	artifactRepo    domain.ArtifactRepository
	appRepo         domain.ApplicationRepository
	adRepo          domain.AvailableDomainRepository
	buildRepo       domain.BuildRepository
	envRepo         domain.EnvironmentRepository
	gitRepo         domain.GitRepositoryRepository
	storage         domain.Storage
	component       domain.ComponentService
	mariaDBManager  domain.MariaDBManager
	mongoDBManager  domain.MongoDBManager
	containerLogger domain.ContainerLogger
	logSvc          *LogStreamService
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
	component domain.ComponentService,
	mariaDBManager domain.MariaDBManager,
	mongoDBManager domain.MongoDBManager,
	containerLogger domain.ContainerLogger,
	logSvc *LogStreamService,
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
		component:       component,
		mariaDBManager:  mariaDBManager,
		mongoDBManager:  mongoDBManager,
		containerLogger: containerLogger,
		logSvc:          logSvc,
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

func (s *APIServerService) GetRepositories(ctx context.Context) ([]*domain.Repository, error) {
	return s.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{})
}

func (s *APIServerService) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	if !repo.IsValid() {
		return newError(ErrorTypeBadRequest, "invalid repository", nil)
	}

	return s.gitRepo.CreateRepository(ctx, repo)
}

func (s *APIServerService) UpdateRepository(ctx context.Context, id string, args *domain.UpdateRepositoryArgs) error {
	err := s.isRepositoryOwner(ctx, id)
	if err != nil {
		return err
	}

	repo, err := s.gitRepo.GetRepository(ctx, id)
	if err != nil {
		return err
	}
	repo.Apply(args)
	if !repo.IsValid() {
		return newError(ErrorTypeBadRequest, "invalid repository", nil)
	}

	return s.gitRepo.UpdateRepository(ctx, id, args)
}

func (s *APIServerService) DeleteRepository(ctx context.Context, id string) error {
	err := s.isRepositoryOwner(ctx, id)
	if err != nil {
		return err
	}

	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{RepositoryID: optional.From(id)})
	if err != nil {
		return errors.Wrap(err, "failed to get related applications")
	}
	if len(apps) > 0 {
		return newError(ErrorTypeBadRequest, "all related applications must be deleted first", nil)
	}

	return s.gitRepo.DeleteRepository(ctx, id)
}

func (s *APIServerService) GetApplications(ctx context.Context) ([]*domain.Application, error) {
	return s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
}

func (s *APIServerService) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	return s.adRepo.GetAvailableDomains(ctx)
}

func (s *APIServerService) AddAvailableDomain(ctx context.Context, ad *domain.AvailableDomain) error {
	err := s.isAdmin(ctx)
	if err != nil {
		return err
	}

	if !ad.IsValid() {
		return newError(ErrorTypeBadRequest, "invalid new domain", nil)
	}
	return s.adRepo.AddAvailableDomain(ctx, ad)
}

func (s *APIServerService) CreateApplication(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	repo, err := s.gitRepo.GetRepository(ctx, app.RepositoryID)
	if err != nil {
		return nil, err
	}
	if repo.Auth.Valid {
		err = s.isRepositoryOwner(ctx, app.RepositoryID)
		if err != nil {
			return nil, err
		}
	}

	if !app.IsValid() {
		return nil, newError(ErrorTypeBadRequest, "invalid application", nil)
	}

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

func (s *APIServerService) createApplicationDatabase(ctx context.Context, app *domain.Application) error {
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

func (s *APIServerService) deleteApplicationDatabase(ctx context.Context, app *domain.Application) error {
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

func (s *APIServerService) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	application, err := s.appRepo.GetApplication(ctx, id)
	return handleRepoError(application, err)
}

func (s *APIServerService) UpdateApplication(ctx context.Context, id string, args *domain.UpdateApplicationArgs) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return err
	}
	app.Apply(args)
	if !app.IsValid() {
		return newError(ErrorTypeBadRequest, "invalid application", nil)
	}

	return s.appRepo.UpdateApplication(ctx, id, args)
}

func (s *APIServerService) DeleteApplication(ctx context.Context, id string) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

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

func (s *APIServerService) GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(applicationID)})
}

func (s *APIServerService) GetBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	return handleRepoError(build, err)
}

func (s *APIServerService) GetBuildLog(ctx context.Context, buildID string) ([]byte, error) {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return nil, err
	}

	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		return nil, err
	}
	if !build.Status.IsFinished() {
		return nil, newError(ErrorTypeBadRequest, "build not finished", nil)
	}
	return domain.GetBuildLog(s.storage, buildID)
}

func (s *APIServerService) GetBuildLogStream(ctx context.Context, buildID string, send func(b []byte) error) error {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return err
	}

	sub := make(chan []byte, 100)
	ok, unsubscribe := s.logSvc.SubscribeBuildLog(buildID, sub)
	if !ok {
		return newError(ErrorTypeBadRequest, "build log stream not available", nil)
	}
	defer unsubscribe()

	for b := range sub {
		err = send(b)
		if err != nil {
			return errors.Wrap(err, "failed to send log")
		}
	}
	return nil
}

func (s *APIServerService) GetArtifact(_ context.Context, artifactID string) ([]byte, error) {
	return domain.GetArtifact(s.storage, artifactID)
}

func (s *APIServerService) GetEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	return s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(applicationID)})
}

func (s *APIServerService) SetEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return err
	}

	env := &domain.Environment{ApplicationID: applicationID, Key: key, Value: value, System: false}
	return s.envRepo.SetEnv(ctx, env)
}

func (s *APIServerService) GetOutput(ctx context.Context, id string, before time.Time) ([]*domain.ContainerLog, error) {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.containerLogger.Get(ctx, id, before)
}

func (s *APIServerService) GetOutputStream(ctx context.Context, id string, after time.Time, send func(l *domain.ContainerLog) error) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := s.containerLogger.Stream(ctx, id, after)
	if err != nil {
		return errors.Wrap(err, "failed to connect to stream")
	}

	for {
		select {
		case d, ok := <-ch:
			if !ok {
				return errors.Wrap(err, "log stream closed")
			}
			err = send(d)
			if err != nil {
				return errors.Wrap(err, "failed to send log")
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *APIServerService) CancelBuild(ctx context.Context, buildID string) error {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return err
	}

	s.component.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_CANCEL_BUILD,
		Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.BuildIdRequest{BuildId: buildID}},
	})
	return nil
}

func (s *APIServerService) RetryCommitBuild(ctx context.Context, applicationID string, commit string) error {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return err
	}

	err = s.buildRepo.MarkCommitAsRetriable(ctx, applicationID, commit)
	if err != nil {
		return err
	}
	// NOTE: requires the app to be running for builds to register
	s.bus.Publish(event.CDServiceRegisterBuildRequest, nil)
	return nil
}

func (s *APIServerService) StartApplication(ctx context.Context, id string) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

	err = s.appRepo.UpdateApplication(ctx, id, &domain.UpdateApplicationArgs{
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

func (s *APIServerService) StopApplication(ctx context.Context, id string) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

	err = s.appRepo.UpdateApplication(ctx, id, &domain.UpdateApplicationArgs{
		Running:   optional.From(false),
		UpdatedAt: optional.From(time.Now()),
	})
	if err != nil {
		return errors.Wrap(err, "failed to mark application as not running")
	}
	s.bus.Publish(event.CDServiceSyncDeployRequest, nil)
	return nil
}
