package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

var BuildStatusMapper = mapper.MustNewValueMapper(map[string]domain.BuildStatus{
	models.BuildsStatusQueued:    domain.BuildStatusQueued,
	models.BuildsStatusBuilding:  domain.BuildStatusBuilding,
	models.BuildsStatusSucceeded: domain.BuildStatusSucceeded,
	models.BuildsStatusFailed:    domain.BuildStatusFailed,
	models.BuildsStatusCanceled:  domain.BuildStatusCanceled,
	models.BuildsStatusSkipped:   domain.BuildStatusSkipped,
})

func FromDomainBuild(build *domain.Build) *models.Build {
	return &models.Build{
		ID:            build.ID,
		Commit:        build.Commit,
		Status:        BuildStatusMapper.FromMust(build.Status),
		QueuedAt:      build.QueuedAt,
		StartedAt:     optional.IntoTime(build.StartedAt),
		UpdatedAt:     optional.IntoTime(build.UpdatedAt),
		FinishedAt:    optional.IntoTime(build.FinishedAt),
		Retriable:     build.Retriable,
		ApplicationID: build.ApplicationID,
	}
}

func ToDomainBuild(build *models.Build) *domain.Build {
	ret := &domain.Build{
		ID:            build.ID,
		Commit:        build.Commit,
		Status:        BuildStatusMapper.IntoMust(build.Status),
		ApplicationID: build.ApplicationID,
		QueuedAt:      build.QueuedAt,
		StartedAt:     optional.FromTime(build.StartedAt),
		UpdatedAt:     optional.FromTime(build.UpdatedAt),
		FinishedAt:    optional.FromTime(build.FinishedAt),
		Retriable:     build.Retriable,
	}
	if build.R != nil && build.R.Artifact != nil {
		ret.Artifact = optional.From(*ToDomainArtifact(build.R.Artifact))
	}
	return ret
}
