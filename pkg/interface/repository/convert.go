package repository

import (
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func fromDomainArtifact(a *domain.Artifact) *models.Artifact {
	return &models.Artifact{
		ID:        a.ID,
		Size:      a.Size,
		CreatedAt: a.CreatedAt,
		DeletedAt: optional.IntoTime(a.DeletedAt),
		BuildID:   a.BuildID,
	}
}

func toDomainArtifact(a *models.Artifact) *domain.Artifact {
	return &domain.Artifact{
		ID:        a.ID,
		BuildID:   a.BuildID,
		Size:      a.Size,
		CreatedAt: a.CreatedAt,
		DeletedAt: optional.FromTime(a.DeletedAt),
	}
}

func fromDomainAvailableDomain(ad *domain.AvailableDomain) *models.AvailableDomain {
	return &models.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}

func toDomainAvailableDomain(ad *models.AvailableDomain) *domain.AvailableDomain {
	return &domain.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}

func fromDomainApplicationConfig(appID string, c *domain.ApplicationConfig) *models.ApplicationConfig {
	return &models.ApplicationConfig{
		ApplicationID:  appID,
		UseMariadb:     c.UseMariaDB,
		UseMongodb:     c.UseMongoDB,
		BaseImage:      c.BaseImage,
		DockerfileName: c.DockerfileName,
		ArtifactPath:   c.ArtifactPath,
		BuildCMD:       c.BuildCmd,
		EntrypointCMD:  c.EntrypointCmd,
		Authentication: c.Authentication.String(),
	}
}

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

func fromDomainRepository(repo *domain.Repository) *models.Repository {
	return &models.Repository{
		ID:   repo.ID,
		Name: repo.Name,
		URL:  repo.URL,
	}
}

func fromDomainRepositoryAuth(repositoryID string, auth *domain.RepositoryAuth) *models.RepositoryAuth {
	return &models.RepositoryAuth{
		RepositoryID: repositoryID,
		Method:       auth.Method.String(),
		Username:     auth.Username,
		Password:     auth.Password,
		SSHKey:       auth.SSHKey,
	}
}

func toDomainRepository(repo *models.Repository) *domain.Repository {
	ret := &domain.Repository{
		ID:       repo.ID,
		Name:     repo.Name,
		URL:      repo.URL,
		OwnerIDs: lo.Map(repo.R.Users, func(user *models.User, i int) string { return user.ID }),
	}
	if repo.R.RepositoryAuth != nil {
		auth := repo.R.RepositoryAuth
		ret.Auth = optional.From(domain.RepositoryAuth{
			Method:   domain.RepositoryAuthMethodFromString(auth.Method),
			Username: auth.Username,
			Password: auth.Password,
			SSHKey:   auth.SSHKey,
		})
	}
	return ret
}

func fromDomainApplication(app *domain.Application) *models.Application {
	return &models.Application{
		ID:            app.ID,
		Name:          app.Name,
		RepositoryID:  app.RepositoryID,
		BranchName:    app.BranchName,
		BuildType:     app.BuildType.String(),
		State:         app.State.String(),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,
	}
}

func toDomainApplication(app *models.Application) *domain.Application {
	return &domain.Application{
		ID:            app.ID,
		Name:          app.Name,
		RepositoryID:  app.RepositoryID,
		BranchName:    app.BranchName,
		BuildType:     builder.BuildTypeFromString(app.BuildType),
		State:         domain.ApplicationStateFromString(app.State),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,

		Config:   toDomainApplicationConfig(app.R.ApplicationConfig),
		Websites: lo.Map(app.R.Websites, func(website *models.Website, i int) *domain.Website { return toDomainWebsite(website) }),
		OwnerIDs: lo.Map(app.R.Users, func(user *models.User, i int) string { return user.ID }),
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
		ret.Artifact = optional.From(*toDomainArtifact(build.R.Artifact))
	}
	return ret
}

func fromDomainEnvironment(applicationID string, env *domain.Environment) *models.Environment {
	return &models.Environment{
		ApplicationID: applicationID,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}

func toDomainEnvironment(env *models.Environment) *domain.Environment {
	return &domain.Environment{
		Key:    env.Key,
		Value:  env.Value,
		System: env.System,
	}
}

func fromDomainWebsite(website *domain.Website) *models.Website {
	return &models.Website{
		ID:          website.ID,
		FQDN:        website.FQDN,
		PathPrefix:  website.PathPrefix,
		StripPrefix: website.StripPrefix,
		HTTPS:       website.HTTPS,
		HTTPPort:    website.HTTPPort,
	}
}

func toDomainWebsite(website *models.Website) *domain.Website {
	return &domain.Website{
		ID:          website.ID,
		FQDN:        website.FQDN,
		PathPrefix:  website.PathPrefix,
		StripPrefix: website.StripPrefix,
		HTTPS:       website.HTTPS,
		HTTPPort:    website.HTTPPort,
	}
}
