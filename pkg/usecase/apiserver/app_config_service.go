package apiserver

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *Service) GetEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	return s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(applicationID)})
}

func (s *Service) SetEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return err
	}

	env := &domain.Environment{ApplicationID: applicationID, Key: key, Value: value, System: false}
	return s.envRepo.SetEnv(ctx, env)
}

func (s *Service) DeleteEnvironmentVariable(ctx context.Context, applicationID string, key string) error {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return err
	}

	return s.envRepo.DeleteEnv(ctx, domain.GetEnvCondition{
		ApplicationID: optional.From(applicationID),
		Key:           optional.From(key),
	})
}

func (s *Service) GetOutputStream(ctx context.Context, id string, send func(l *domain.ContainerLog) error) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := s.containerLogger.Stream(ctx, id)
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

func (s *Service) StartApplication(ctx context.Context, id string) error {
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

	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return errors.Wrap(err, "getting application")
	}
	err = s.controller.FetchRepository(ctx, app.RepositoryID)
	if err != nil {
		return errors.Wrap(err, "failed to request repository fetch")
	}
	err = s.controller.SyncDeployments(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to request sync deployment")
	}
	return nil
}

func (s *Service) StopApplication(ctx context.Context, id string) error {
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

	err = s.controller.SyncDeployments(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to request sync deployment")
	}
	return nil
}
