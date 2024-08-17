package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
)

func FromDomainRuntimeImage(r *domain.RuntimeImage) *models.RuntimeImage {
	return &models.RuntimeImage{
		ID:      r.ID,
		BuildID: r.BuildID,
		Size:    r.Size,
	}
}

func ToDomainRuntimeImage(r *models.RuntimeImage) domain.RuntimeImage {
	return domain.RuntimeImage{
		ID:        r.ID,
		BuildID:   r.BuildID,
		Size:      r.Size,
		CreatedAt: r.CreatedAt,
	}
}
