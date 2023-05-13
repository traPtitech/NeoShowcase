package usecase

import (
	"context"

	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIServerService) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	return s.adRepo.GetAvailableDomains(ctx)
}

func (s *APIServerService) AddAvailableDomain(ctx context.Context, ad *domain.AvailableDomain) error {
	err := s.isAdmin(ctx)
	if err != nil {
		return err
	}

	if err = ad.Validate(); err != nil {
		return newError(ErrorTypeBadRequest, "invalid new domain", err)
	}
	return s.adRepo.AddAvailableDomain(ctx, ad)
}

func (s *APIServerService) GetAvailablePorts(ctx context.Context) (domain.AvailablePortSlice, []*domain.UnavailablePort, error) {
	available, err := s.apRepo.GetAvailablePorts(ctx)
	if err != nil {
		return nil, nil, err
	}
	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return nil, nil, err
	}
	usedPorts := ds.Map(apps, func(app *domain.Application) []*domain.PortPublication { return app.PortPublications })
	unavailable := ds.Map(lo.Flatten(usedPorts), (*domain.PortPublication).ToUnavailablePort)
	return available, unavailable, nil
}

func (s *APIServerService) AddAvailablePort(ctx context.Context, ap *domain.AvailablePort) error {
	err := s.isAdmin(ctx)
	if err != nil {
		return err
	}

	err = ap.Validate()
	if err != nil {
		return newError(ErrorTypeBadRequest, "invalid port definition", err)
	}

	return s.apRepo.AddAvailablePort(ctx, ap)
}
