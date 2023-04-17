package usecase

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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
