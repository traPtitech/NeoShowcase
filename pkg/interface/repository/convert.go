package repository

import (
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func toDomainApplicationConfig(c *models.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		UseMariaDB:     c.UseMariadb,
		UseMongoDB:     c.UseMongodb,
		BaseImage:      c.BaseImage,
		DockerfileName: c.DockerfileName,
		ArtifactPath:   c.ArtifactPath,
		BuildCmd:       c.BuildCMD,
		EntrypointCmd:  c.EntrypointCMD,
		Authentication: domain.AuthenticationTypeFromString(c.Authentication),
	}
}

func toDomainRepository(repo *models.Repository) domain.Repository {
	return domain.Repository{
		ID:  repo.ID,
		URL: repo.URL,
	}
}

func toDomainApplication(app *models.Application) *domain.Application {
	return &domain.Application{
		ID:            app.ID,
		Name:          app.Name,
		BranchName:    app.BranchName,
		BuildType:     builder.BuildTypeFromString(app.BuildType),
		State:         domain.ApplicationStateFromString(app.State),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,

		Config:     toDomainApplicationConfig(app.R.ApplicationConfig),
		Repository: toDomainRepository(app.R.Repository),
		Websites:   lo.Map(app.R.Websites, func(website *models.Website, i int) *domain.Website { return toDomainWebsite(website) }),
	}
}

func fromDomainBuild(build *domain.Build) *models.Build {
	return &models.Build{
		ID:            build.ID,
		Commit:        build.Commit,
		Status:        build.Status.String(),
		StartedAt:     optional.IntoTime(build.StartedAt),
		UpdatedAt:     optional.IntoTime(build.UpdatedAt),
		FinishedAt:    optional.IntoTime(build.FinishedAt),
		Retriable:     build.Retriable,
		ApplicationID: build.ApplicationID,
	}
}

func toDomainBuild(build *models.Build) *domain.Build {
	ret := &domain.Build{
		ID:            build.ID,
		Commit:        build.Commit,
		Status:        builder.BuildStatusFromString(build.Status),
		ApplicationID: build.ApplicationID,
		StartedAt:     optional.FromTime(build.StartedAt),
		UpdatedAt:     optional.FromTime(build.UpdatedAt),
		FinishedAt:    optional.FromTime(build.FinishedAt),
		Retriable:     build.Retriable,
	}
	if build.R != nil && build.R.Artifact != nil {
		artifact := build.R.Artifact
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

func fromDomainWebsite(website *domain.Website) *models.Website {
	return &models.Website{
		ID:         website.ID,
		FQDN:       website.FQDN,
		PathPrefix: website.PathPrefix,
		HTTPS:      website.HTTPS,
		HTTPPort:   website.HTTPPort,
	}
}

func toDomainWebsite(website *models.Website) *domain.Website {
	return &domain.Website{
		ID:         website.ID,
		FQDN:       website.FQDN,
		PathPrefix: website.PathPrefix,
		HTTPS:      website.HTTPS,
		HTTPPort:   website.HTTPPort,
	}
}
