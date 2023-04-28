package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var BuildTypeMapper = mapper.MustNewValueMapper(map[string]domain.BuildType{
	models.ApplicationConfigBuildTypeRuntimeBuildpack:  domain.BuildTypeRuntimeBuildpack,
	models.ApplicationConfigBuildTypeRuntimeCMD:        domain.BuildTypeRuntimeCmd,
	models.ApplicationConfigBuildTypeRuntimeDockerfile: domain.BuildTypeRuntimeDockerfile,
	models.ApplicationConfigBuildTypeStaticCMD:         domain.BuildTypeStaticCmd,
	models.ApplicationConfigBuildTypeStaticDockerfile:  domain.BuildTypeStaticDockerfile,
})

func FromDomainBuildConfig(c domain.BuildConfig, mc *models.ApplicationConfig) {
	switch bc := c.(type) {
	case *domain.BuildConfigRuntimeBuildpack:
		mc.Context = bc.Context
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
}

func ToDomainBuildConfig(c *models.ApplicationConfig) domain.BuildConfig {
	switch BuildTypeMapper.IntoMust(c.BuildType) {
	case domain.BuildTypeRuntimeBuildpack:
		return &domain.BuildConfigRuntimeBuildpack{
			Context: c.Context,
		}
	case domain.BuildTypeRuntimeCmd:
		return &domain.BuildConfigRuntimeCmd{
			BaseImage:     c.BaseImage,
			BuildCmd:      c.BuildCMD,
			BuildCmdShell: c.BuildCMDShell,
		}
	case domain.BuildTypeRuntimeDockerfile:
		return &domain.BuildConfigRuntimeDockerfile{
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
