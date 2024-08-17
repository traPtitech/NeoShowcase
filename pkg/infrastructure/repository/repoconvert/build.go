package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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
		ConfigHash:    build.ConfigHash,
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
		ConfigHash:    build.ConfigHash,
		Status:        BuildStatusMapper.IntoMust(build.Status),
		ApplicationID: build.ApplicationID,
		QueuedAt:      build.QueuedAt,
		StartedAt:     optional.FromTime(build.StartedAt),
		UpdatedAt:     optional.FromTime(build.UpdatedAt),
		FinishedAt:    optional.FromTime(build.FinishedAt),
		Retriable:     build.Retriable,
	}
	if build.R != nil {
		ret.Artifacts = ds.Map(build.R.Artifacts, ToDomainArtifact)
		if len(build.R.RuntimeImages) > 0 {
			ret.RuntimeImage = optional.From(ToDomainRuntimeImage(build.R.RuntimeImages[0]))
		}
	}
	return ret
}
