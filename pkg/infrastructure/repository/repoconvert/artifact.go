package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func FromDomainArtifact(a *domain.Artifact) *models.Artifact {
	return &models.Artifact{
		ID:        a.ID,
		Name:      a.Name,
		Size:      a.Size,
		CreatedAt: a.CreatedAt,
		DeletedAt: optional.IntoTime(a.DeletedAt),
		BuildID:   a.BuildID,
	}
}

func ToDomainArtifact(a *models.Artifact) *domain.Artifact {
	return &domain.Artifact{
		ID:        a.ID,
		Name:      a.Name,
		BuildID:   a.BuildID,
		Size:      a.Size,
		CreatedAt: a.CreatedAt,
		DeletedAt: optional.FromTime(a.DeletedAt),
	}
}
