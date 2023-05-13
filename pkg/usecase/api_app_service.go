package usecase

import (
	"context"
	"strconv"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/random"
)

func (s *APIServerService) CreateApplication(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	repo, err := s.gitRepo.GetRepository(ctx, app.RepositoryID)
	if err != nil {
		return nil, err
	}

	// Only check for repository owner if repository is private;
	// allow everyone to create application if repository is public
	if repo.Auth.Valid {
		err = s.isRepositoryOwner(ctx, app.RepositoryID)
		if err != nil {
			return nil, err
		}
	}

	// Validate
	existingApps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return nil, errors.Wrap(err, "getting existing applications")
	}
	domains, err := s.adRepo.GetAvailableDomains(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting available domains")
	}
	ports, err := s.apRepo.GetAvailablePorts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting available ports")
	}
	valErr, err := app.Validate(ctx, web.GetUser(ctx), s.controller, existingApps, domains, ports)
	if err != nil {
		return nil, errors.Wrap(err, "validating application")
	}
	if valErr != nil {
		return nil, newError(ErrorTypeBadRequest, "invalid application", valErr)
	}

	// Create
	err = s.appRepo.CreateApplication(ctx, app)
	if err != nil {
		return nil, err
	}

	err = s.createApplicationDatabase(ctx, app)
	if err != nil {
		return nil, err
	}

	err = s.controller.FetchRepository(ctx, app.RepositoryID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request to fetch repository")
	}

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

func (s *APIServerService) GetApplications(ctx context.Context) ([]*domain.Application, error) {
	return s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
}

func (s *APIServerService) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	return handleRepoError(s.appRepo.GetApplication(ctx, id))
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

	// Validate
	existingApps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return errors.Wrap(err, "getting existing applications")
	}
	domains, err := s.adRepo.GetAvailableDomains(ctx)
	if err != nil {
		return errors.Wrap(err, "getting available domains")
	}
	ports, err := s.apRepo.GetAvailablePorts(ctx)
	if err != nil {
		return errors.Wrap(err, "getting available ports")
	}
	valErr, err := app.Validate(ctx, web.GetUser(ctx), s.controller, existingApps, domains, ports)
	if err != nil {
		return errors.Wrap(err, "validating application")
	}
	if valErr != nil {
		return newError(ErrorTypeBadRequest, "invalid application", valErr)
	}

	// Update
	err = s.appRepo.UpdateApplication(ctx, id, args)
	if err != nil {
		return errors.Wrap(err, "updating application")
	}

	// Sync
	err = s.controller.FetchRepository(ctx, app.RepositoryID)
	if err != nil {
		return errors.Wrap(err, "requesting fetch repository")
	}
	err = s.controller.SyncDeployments(ctx)
	if err != nil {
		return errors.Wrap(err, "requesting sync deployments")
	}

	return nil
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
