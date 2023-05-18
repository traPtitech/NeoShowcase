package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func FromPBRuntimeConfig(c *pb.RuntimeConfig) domain.RuntimeConfig {
	return domain.RuntimeConfig{
		UseMariaDB: c.UseMariadb,
		UseMongoDB: c.UseMongodb,
		Entrypoint: c.Entrypoint,
		Command:    c.Command,
	}
}

func ToPBRuntimeConfig(c *domain.RuntimeConfig) *pb.RuntimeConfig {
	return &pb.RuntimeConfig{
		UseMariadb: c.UseMariaDB,
		UseMongodb: c.UseMongoDB,
		Entrypoint: c.Entrypoint,
		Command:    c.Command,
	}
}

func FromPBBuildConfig(c *pb.ApplicationConfig) domain.BuildConfig {
	switch bc := c.BuildConfig.(type) {
	case *pb.ApplicationConfig_RuntimeBuildpack:
		return &domain.BuildConfigRuntimeBuildpack{
			RuntimeConfig: FromPBRuntimeConfig(bc.RuntimeBuildpack.RuntimeConfig),
			Context:       bc.RuntimeBuildpack.Context,
		}
	case *pb.ApplicationConfig_RuntimeCmd:
		return &domain.BuildConfigRuntimeCmd{
			RuntimeConfig: FromPBRuntimeConfig(bc.RuntimeCmd.RuntimeConfig),
			BaseImage:     bc.RuntimeCmd.BaseImage,
			BuildCmd:      bc.RuntimeCmd.BuildCmd,
			BuildCmdShell: bc.RuntimeCmd.BuildCmdShell,
		}
	case *pb.ApplicationConfig_RuntimeDockerfile:
		return &domain.BuildConfigRuntimeDockerfile{
			RuntimeConfig:  FromPBRuntimeConfig(bc.RuntimeDockerfile.RuntimeConfig),
			DockerfileName: bc.RuntimeDockerfile.DockerfileName,
			Context:        bc.RuntimeDockerfile.Context,
		}
	case *pb.ApplicationConfig_StaticCmd:
		return &domain.BuildConfigStaticCmd{
			BaseImage:     bc.StaticCmd.BaseImage,
			BuildCmd:      bc.StaticCmd.BuildCmd,
			BuildCmdShell: bc.StaticCmd.BuildCmdShell,
			ArtifactPath:  bc.StaticCmd.ArtifactPath,
		}
	case *pb.ApplicationConfig_StaticDockerfile:
		return &domain.BuildConfigStaticDockerfile{
			DockerfileName: bc.StaticDockerfile.DockerfileName,
			Context:        bc.StaticDockerfile.Context,
			ArtifactPath:   bc.StaticDockerfile.ArtifactPath,
		}
	default:
		panic("unknown pb build config type")
	}
}

func FromPBApplicationConfig(c *pb.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		BuildConfig: FromPBBuildConfig(c),
	}
}

func ToPBApplicationConfig(c domain.ApplicationConfig) *pb.ApplicationConfig {
	switch bc := c.BuildConfig.(type) {
	case *domain.BuildConfigRuntimeBuildpack:
		return &pb.ApplicationConfig{
			BuildConfig: &pb.ApplicationConfig_RuntimeBuildpack{RuntimeBuildpack: &pb.BuildConfigRuntimeBuildpack{
				RuntimeConfig: ToPBRuntimeConfig(&bc.RuntimeConfig),
				Context:       bc.Context,
			}},
		}
	case *domain.BuildConfigRuntimeCmd:
		return &pb.ApplicationConfig{
			BuildConfig: &pb.ApplicationConfig_RuntimeCmd{RuntimeCmd: &pb.BuildConfigRuntimeCmd{
				RuntimeConfig: ToPBRuntimeConfig(&bc.RuntimeConfig),
				BaseImage:     bc.BaseImage,
				BuildCmd:      bc.BuildCmd,
				BuildCmdShell: bc.BuildCmdShell,
			}},
		}
	case *domain.BuildConfigRuntimeDockerfile:
		return &pb.ApplicationConfig{
			BuildConfig: &pb.ApplicationConfig_RuntimeDockerfile{RuntimeDockerfile: &pb.BuildConfigRuntimeDockerfile{
				RuntimeConfig:  ToPBRuntimeConfig(&bc.RuntimeConfig),
				DockerfileName: bc.DockerfileName,
				Context:        bc.Context,
			}},
		}
	case *domain.BuildConfigStaticCmd:
		return &pb.ApplicationConfig{
			BuildConfig: &pb.ApplicationConfig_StaticCmd{StaticCmd: &pb.BuildConfigStaticCmd{
				BaseImage:     bc.BaseImage,
				BuildCmd:      bc.BuildCmd,
				BuildCmdShell: bc.BuildCmdShell,
				ArtifactPath:  bc.ArtifactPath,
			}},
		}
	case *domain.BuildConfigStaticDockerfile:
		return &pb.ApplicationConfig{
			BuildConfig: &pb.ApplicationConfig_StaticDockerfile{StaticDockerfile: &pb.BuildConfigStaticDockerfile{
				DockerfileName: bc.DockerfileName,
				Context:        bc.Context,
				ArtifactPath:   bc.ArtifactPath,
			}},
		}
	default:
		panic("unknown domain build config type")
	}
}
