package apiserver

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (s *Service) GetAvailableMetrics(_ context.Context) []string {
	return s.metricsService.AvailableNames()
}

func (s *Service) GetApplicationMetrics(ctx context.Context, name string, id string, before time.Time, limit time.Duration) ([]*domain.AppMetric, error) {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return nil, err
	}
	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.metricsService.Get(ctx, name, app, before, limit)
}

func (s *Service) GetOutput(ctx context.Context, id string, before time.Time, limit int) ([]*domain.ContainerLog, error) {
	// Validate
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return nil, err
	}
	if limit > s.containerLogger.LogLimit() {
		return nil, newError(ErrorTypeBadRequest, "limit too large", nil)
	}

	// Get logs
	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.containerLogger.Get(ctx, app, before, limit)
}

func (s *Service) GetOutputStream(ctx context.Context, id string, begin time.Time, send func(l *domain.ContainerLog) error) error {
	err := s.isApplicationOwner(ctx, id)
	if err != nil {
		return err
	}

	app, err := s.appRepo.GetApplication(ctx, id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := s.containerLogger.Stream(ctx, app, begin)
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
