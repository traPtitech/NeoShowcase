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
	models.ApplicationConfigBuildTypeStaticCMD:         domain.BuildTypeStaticCmd,
	models.ApplicationConfigBuildTypeStaticDockerfile:  domain.BuildTypeStaticDockerfile,
})

func assignRuntimeConfig(mc *models.ApplicationConfig, c *domain.RuntimeConfig) {
	mc.UseMariadb = c.UseMariaDB
	mc.UseMongodb = c.UseMongoDB
	mc.Entrypoint = c.Entrypoint
	mc.Command = c.Command
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
		mc.BuildCMDShell = bc.BuildCmdShell
	case *domain.BuildConfigRuntimeDockerfile:
		assignRuntimeConfig(mc, &bc.RuntimeConfig)
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
}

func ToDomainRuntimeConfig(c *models.ApplicationConfig) domain.RuntimeConfig {
	return domain.RuntimeConfig{
		UseMariaDB: c.UseMariadb,
		UseMongoDB: c.UseMongodb,
		Entrypoint: c.Entrypoint,
		Command:    c.Command,
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
			BuildCmdShell: c.BuildCMDShell,
		}
	case domain.BuildTypeRuntimeDockerfile:
		return &domain.BuildConfigRuntimeDockerfile{
			RuntimeConfig:  ToDomainRuntimeConfig(c),
			DockerfileName: c.DockerfileName,
			Context:        c.Context,
		}
	case domain.BuildTypeStaticCmd:
		return &domain.BuildConfigStaticCmd{
			BaseImage:     c.BaseImage,
			BuildCmd:      c.BuildCMD,
			BuildCmdShell: c.BuildCMDShell,
			ArtifactPath:  c.ArtifactPath,
		}
	case domain.BuildTypeStaticDockerfile:
		return &domain.BuildConfigStaticDockerfile{
			DockerfileName: c.DockerfileName,
			Context:        c.Context,
			ArtifactPath:   c.ArtifactPath,
		}
	default:
		panic("unknown build type")
	}
}
