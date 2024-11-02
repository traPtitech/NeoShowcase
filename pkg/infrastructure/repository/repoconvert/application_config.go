package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

func FromDomainApplicationConfig(appID string, c *domain.ApplicationConfig) *models.ApplicationConfig {
	mc := &models.ApplicationConfig{
		ApplicationID: appID,
		BuildType:     BuildTypeMapper.FromMust(c.BuildConfig.BuildType()),
	}
	assignBuildConfig(mc, c.BuildConfig)
	return mc
}

func ToDomainApplicationConfig(c *models.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		BuildConfig: ToDomainBuildConfig(c),
	}
}

var BuildTypeMapper = mapper.MustNewValueMapper(map[string]domain.BuildType{
	models.ApplicationConfigBuildTypeRuntimeBuildpack:  domain.BuildTypeRuntimeBuildpack,
	models.ApplicationConfigBuildTypeRuntimeCMD:        domain.BuildTypeRuntimeCmd,
	models.ApplicationConfigBuildTypeRuntimeDockerfile: domain.BuildTypeRuntimeDockerfile,
	models.ApplicationConfigBuildTypeStaticBuildpack:   domain.BuildTypeStaticBuildpack,
	models.ApplicationConfigBuildTypeStaticCMD:         domain.BuildTypeStaticCmd,
	models.ApplicationConfigBuildTypeStaticDockerfile:  domain.BuildTypeStaticDockerfile,
})

var StartupBehaviorMapper = mapper.MustNewValueMapper(map[string]domain.StartupBehavior{
	models.ApplicationConfigStartupBehaviorUndefined:   domain.StartupBehaviorUndefined,
	models.ApplicationConfigStartupBehaviorLoadingPage: domain.StartupBehaviorLoadingPage,
	models.ApplicationConfigStartupBehaviorBlocking:    domain.StartupBehaviorBlocking,
})

func assignRuntimeConfig(mc *models.ApplicationConfig, c *domain.RuntimeConfig) {
	mc.UseMariadb = c.UseMariaDB
	mc.UseMongodb = c.UseMongoDB
	mc.AutoShutdown = c.AutoShutdown.Enabled
	mc.StartupBehavior = StartupBehaviorMapper.FromMust(c.AutoShutdown.Startup)
	mc.Entrypoint = c.Entrypoint
	mc.Command = c.Command
}

func ToDomainRuntimeConfig(c *models.ApplicationConfig) domain.RuntimeConfig {
	return domain.RuntimeConfig{
		UseMariaDB: c.UseMariadb,
		UseMongoDB: c.UseMongodb,
		AutoShutdown: domain.AutoShutdownConfig{
			Enabled: c.AutoShutdown,
			Startup: StartupBehaviorMapper.IntoMust(c.StartupBehavior),
		},
		Entrypoint: c.Entrypoint,
		Command:    c.Command,
	}
}

func assignStaticConfig(mc *models.ApplicationConfig, c *domain.StaticConfig) {
	mc.ArtifactPath = c.ArtifactPath
	mc.Spa = c.SPA
}

func ToDomainStaticConfig(c *models.ApplicationConfig) domain.StaticConfig {
	return domain.StaticConfig{
		ArtifactPath: c.ArtifactPath,
		SPA:          c.Spa,
	}
}

func assignBuildConfig(mc *models.ApplicationConfig, c domain.BuildConfig) {
	switch bc := c.(type) {
	case *domain.BuildConfigRuntimeBuildpack:
		assignRuntimeConfig(mc, &bc.RuntimeConfig)
		mc.Context = bc.Context
	case *domain.BuildConfigRuntimeCmd:
		assignRuntimeConfig(mc, &bc.RuntimeConfig)
		mc.BaseImage = bc.BaseImage
		mc.BuildCMD = bc.BuildCmd
	case *domain.BuildConfigRuntimeDockerfile:
		assignRuntimeConfig(mc, &bc.RuntimeConfig)
		mc.DockerfileName = bc.DockerfileName
		mc.Context = bc.Context
	case *domain.BuildConfigStaticBuildpack:
		assignStaticConfig(mc, &bc.StaticConfig)
		mc.Context = bc.Context
	case *domain.BuildConfigStaticCmd:
		assignStaticConfig(mc, &bc.StaticConfig)
		mc.BaseImage = bc.BaseImage
		mc.BuildCMD = bc.BuildCmd
	case *domain.BuildConfigStaticDockerfile:
		assignStaticConfig(mc, &bc.StaticConfig)
		mc.DockerfileName = bc.DockerfileName
		mc.Context = bc.Context
	default:
		panic("unknown domain build config type")
	}
}

func ToDomainBuildConfig(c *models.ApplicationConfig) domain.BuildConfig {
	switch BuildTypeMapper.IntoMust(c.BuildType) {
	case domain.BuildTypeRuntimeBuildpack:
		return &domain.BuildConfigRuntimeBuildpack{
			RuntimeConfig: ToDomainRuntimeConfig(c),
			Context:       c.Context,
		}
	case domain.BuildTypeRuntimeCmd:
		return &domain.BuildConfigRuntimeCmd{
			RuntimeConfig: ToDomainRuntimeConfig(c),
			BaseImage:     c.BaseImage,
			BuildCmd:      c.BuildCMD,
		}
	case domain.BuildTypeRuntimeDockerfile:
		return &domain.BuildConfigRuntimeDockerfile{
			RuntimeConfig:  ToDomainRuntimeConfig(c),
			DockerfileName: c.DockerfileName,
			Context:        c.Context,
		}
	case domain.BuildTypeStaticBuildpack:
		return &domain.BuildConfigStaticBuildpack{
			StaticConfig: ToDomainStaticConfig(c),
			Context:      c.Context,
		}
	case domain.BuildTypeStaticCmd:
		return &domain.BuildConfigStaticCmd{
			StaticConfig: ToDomainStaticConfig(c),
			BaseImage:    c.BaseImage,
			BuildCmd:     c.BuildCMD,
		}
	case domain.BuildTypeStaticDockerfile:
		return &domain.BuildConfigStaticDockerfile{
			StaticConfig:   ToDomainStaticConfig(c),
			DockerfileName: c.DockerfileName,
			Context:        c.Context,
		}
	default:
		panic("unknown build type")
	}
}
