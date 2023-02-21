package repository

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
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
	}
}
