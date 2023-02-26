package repository

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func toDomainRepository(repo *models.Repository) domain.Repository {
	return domain.Repository{
		ID:  repo.ID,
		URL: repo.URL,
	}
}

func toDomainApplication(app *models.Application, repo domain.Repository) *domain.Application {
	return &domain.Application{
		ID:         app.ID,
		Repository: repo,
		BranchName: app.BranchName,
		BuildType:  builder.BuildTypeFromString(app.BuildType),
	}
}

func toDomainBuild(build *models.Build) *domain.Build {
	return &domain.Build{
		ID:            build.ID,
		Status:        builder.BuildStatusFromString(build.Status),
		ApplicationID: build.ApplicationID,
		StartedAt:     build.StartedAt,
		FinishedAt:    optional.New(build.FinishedAt.Time, build.FinishedAt.Valid),
	}
}

func toDomainEnvironment(env *models.Environment) *domain.Environment {
	return &domain.Environment{
		ID:            env.ID,
		ApplicationID: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
	}
}
