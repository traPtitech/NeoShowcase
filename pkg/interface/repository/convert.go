package repository

import (
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
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

var buildTypeMapper = mapper.NewValueMapper(map[string]domain.BuildType{
	models.ApplicationConfigBuildTypeRuntimeCMD:        domain.BuildTypeRuntimeCmd,
	models.ApplicationConfigBuildTypeRuntimeDockerfile: domain.BuildTypeRuntimeDockerfile,
	models.ApplicationConfigBuildTypeStaticCMD:         domain.BuildTypeStaticCmd,
	models.ApplicationConfigBuildTypeStaticDockerfile:  domain.BuildTypeStaticDockerfile,
})

func fromDomainApplicationConfig(appID string, c *domain.ApplicationConfig) *models.ApplicationConfig {
	mc := &models.ApplicationConfig{
		ApplicationID: appID,
		UseMariadb:    c.UseMariaDB,
		UseMongodb:    c.UseMongoDB,
		BuildType:     buildTypeMapper.FromMust(c.BuildType),
		Entrypoint:    c.Entrypoint,
		Command:       c.Command,
	}
	switch bc := c.BuildConfig.(type) {
	case *domain.BuildConfigRuntimeCmd:
		mc.BaseImage = bc.BaseImage
		mc.BuildCMD = bc.BuildCmd
		mc.BuildCMDShell = bc.BuildCmdShell
	case *domain.BuildConfigRuntimeDockerfile:
		mc.DockerfileName = bc.DockerfileName
		mc.Context = bc.Context
	case *domain.BuildConfigStaticCmd:
		mc.BaseImage = bc.BaseImage
		mc.BuildCMD = bc.BuildCmd
		mc.BuildCMDShell = bc.BuildCmdShell
		mc.ArtifactPath = bc.ArtifactPath
	case *domain.BuildConfigStaticDockerfile:
		mc.DockerfileName = bc.DockerfileName
		mc.Context = bc.Context
		mc.ArtifactPath = bc.ArtifactPath
	default:
		panic("unknown domain build config type")
	}
	return mc
}

func toDomainApplicationConfig(c *models.ApplicationConfig) domain.ApplicationConfig {
	ac := domain.ApplicationConfig{
		UseMariaDB: c.UseMariadb,
		UseMongoDB: c.UseMongodb,
		BuildType:  buildTypeMapper.IntoMust(c.BuildType),
		Entrypoint: c.Entrypoint,
		Command:    c.Command,
	}
	switch ac.BuildType {
	case domain.BuildTypeRuntimeCmd:
		ac.BuildConfig = &domain.BuildConfigRuntimeCmd{
			BaseImage:     c.BaseImage,
			BuildCmd:      c.BuildCMD,
			BuildCmdShell: c.BuildCMDShell,
		}
	case domain.BuildTypeRuntimeDockerfile:
		ac.BuildConfig = &domain.BuildConfigRuntimeDockerfile{
			DockerfileName: c.DockerfileName,
			Context:        c.Context,
		}
	case domain.BuildTypeStaticCmd:
		ac.BuildConfig = &domain.BuildConfigStaticCmd{
			BaseImage:     c.BaseImage,
			BuildCmd:      c.BuildCMD,
			BuildCmdShell: c.BuildCMDShell,
			ArtifactPath:  c.ArtifactPath,
		}
	case domain.BuildTypeStaticDockerfile:
		ac.BuildConfig = &domain.BuildConfigStaticDockerfile{
			DockerfileName: c.DockerfileName,
			Context:        c.Context,
			ArtifactPath:   c.ArtifactPath,
		}
	default:
		panic("unknown build type")
	}
	return ac
}

func fromDomainRepository(repo *domain.Repository) *models.Repository {
	return &models.Repository{
		ID:   repo.ID,
		Name: repo.Name,
		URL:  repo.URL,
	}
}

var repoAuthMethodMapper = mapper.NewValueMapper(map[string]domain.RepositoryAuthMethod{
	models.RepositoryAuthMethodBasic: domain.RepositoryAuthMethodBasic,
	models.RepositoryAuthMethodSSH:   domain.RepositoryAuthMethodSSH,
})

func fromDomainRepositoryAuth(repositoryID string, auth *domain.RepositoryAuth) *models.RepositoryAuth {
	return &models.RepositoryAuth{
		RepositoryID: repositoryID,
		Method:       repoAuthMethodMapper.FromMust(auth.Method),
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
			Method:   repoAuthMethodMapper.IntoMust(auth.Method),
			Username: auth.Username,
			Password: auth.Password,
			SSHKey:   auth.SSHKey,
		})
	}
	return ret
}

var deployTypeMapper = mapper.NewValueMapper(map[string]domain.DeployType{
	models.ApplicationsDeployTypeRuntime: domain.DeployTypeRuntime,
	models.ApplicationsDeployTypeStatic:  domain.DeployTypeStatic,
})

var containerStateMapper = mapper.NewValueMapper(map[string]domain.ContainerState{
	models.ApplicationsContainerMissing:  domain.ContainerStateMissing,
	models.ApplicationsContainerStarting: domain.ContainerStateStarting,
	models.ApplicationsContainerRunning:  domain.ContainerStateRunning,
	models.ApplicationsContainerExited:   domain.ContainerStateExited,
	models.ApplicationsContainerErrored:  domain.ContainerStateErrored,
	models.ApplicationsContainerUnknown:  domain.ContainerStateUnknown,
})

func fromDomainApplication(app *domain.Application) *models.Application {
	return &models.Application{
		ID:            app.ID,
		Name:          app.Name,
		RepositoryID:  app.RepositoryID,
		RefName:       app.RefName,
		DeployType:    deployTypeMapper.FromMust(app.DeployType),
		Running:       app.Running,
		Container:     containerStateMapper.FromMust(app.Container),
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
		RefName:       app.RefName,
		DeployType:    deployTypeMapper.IntoMust(app.DeployType),
		Running:       app.Running,
		Container:     containerStateMapper.IntoMust(app.Container),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,

		Config:   toDomainApplicationConfig(app.R.ApplicationConfig),
		Websites: lo.Map(app.R.Websites, func(website *models.Website, i int) *domain.Website { return toDomainWebsite(website) }),
		OwnerIDs: lo.Map(app.R.Users, func(user *models.User, i int) string { return user.ID }),
	}
}

var buildStatusMapper = mapper.NewValueMapper(map[string]domain.BuildStatus{
	models.BuildsStatusQueued:    domain.BuildStatusQueued,
	models.BuildsStatusBuilding:  domain.BuildStatusBuilding,
	models.BuildsStatusSucceeded: domain.BuildStatusSucceeded,
	models.BuildsStatusFailed:    domain.BuildStatusFailed,
	models.BuildsStatusCanceled:  domain.BuildStatusCanceled,
	models.BuildsStatusSkipped:   domain.BuildStatusSkipped,
})

func fromDomainBuild(build *domain.Build) *models.Build {
	return &models.Build{
		ID:            build.ID,
		Commit:        build.Commit,
		Status:        buildStatusMapper.FromMust(build.Status),
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
		Status:        buildStatusMapper.IntoMust(build.Status),
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

func fromDomainEnvironment(env *domain.Environment) *models.Environment {
	return &models.Environment{
		ApplicationID: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}

func toDomainEnvironment(env *models.Environment) *domain.Environment {
	return &domain.Environment{
		ApplicationID: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}

var authTypeMapper = mapper.NewValueMapper(map[string]domain.AuthenticationType{
	models.WebsitesAuthenticationOff:  domain.AuthenticationTypeOff,
	models.WebsitesAuthenticationSoft: domain.AuthenticationTypeSoft,
	models.WebsitesAuthenticationHard: domain.AuthenticationTypeHard,
})

func fromDomainWebsite(appID string, website *domain.Website) *models.Website {
	return &models.Website{
		ID:             website.ID,
		FQDN:           website.FQDN,
		PathPrefix:     website.PathPrefix,
		StripPrefix:    website.StripPrefix,
		HTTPS:          website.HTTPS,
		HTTPPort:       website.HTTPPort,
		Authentication: authTypeMapper.FromMust(website.Authentication),
		ApplicationID:  appID,
	}
}

func toDomainWebsite(website *models.Website) *domain.Website {
	return &domain.Website{
		ID:             website.ID,
		FQDN:           website.FQDN,
		PathPrefix:     website.PathPrefix,
		StripPrefix:    website.StripPrefix,
		HTTPS:          website.HTTPS,
		HTTPPort:       website.HTTPPort,
		Authentication: authTypeMapper.IntoMust(website.Authentication),
	}
}

func fromDomainUser(user *domain.User) *models.User {
	return &models.User{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}

func toDomainUser(user *models.User) *domain.User {
	return &domain.User{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}
