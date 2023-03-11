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
		ID:            app.ID,
		Repository:    repo,
		BranchName:    app.BranchName,
		BuildType:     builder.BuildTypeFromString(app.BuildType),
		State:         domain.ApplicationStateFromString(app.State),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
	}
}

func toDomainBuild(build *models.Build) *domain.Build {
	ret := &domain.Build{
		ID:            build.ID,
		Commit:        build.Commit,
		Status:        builder.BuildStatusFromString(build.Status),
		ApplicationID: build.ApplicationID,
		StartedAt:     build.StartedAt,
		FinishedAt:    optional.New(build.FinishedAt.Time, build.FinishedAt.Valid),
	}
	artifact := build.R.Artifact
	if artifact != nil {
		ret.Artifact = optional.From(domain.Artifact{
			ID:        artifact.ID,
			Size:      artifact.Size,
			CreatedAt: artifact.CreatedAt,
		})
	}
	return ret
}

func toDomainEnvironment(env *models.Environment) *domain.Environment {
	return &domain.Environment{
		ID:            env.ID,
		ApplicationID: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
	}
}
