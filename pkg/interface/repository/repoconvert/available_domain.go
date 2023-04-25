package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

func FromDomainAvailableDomain(ad *domain.AvailableDomain) *models.AvailableDomain {
	return &models.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}

func ToDomainAvailableDomain(ad *models.AvailableDomain) *domain.AvailableDomain {
	return &domain.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}
