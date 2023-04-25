package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
)

func FromPBBuildConfig(c *pb.BuildConfig) domain.BuildConfig {
	switch bc := c.BuildConfig.(type) {
	case *pb.BuildConfig_RuntimeCmd:
		return &domain.BuildConfigRuntimeCmd{
			BaseImage:     bc.RuntimeCmd.BaseImage,
			BuildCmd:      bc.RuntimeCmd.BuildCmd,
			BuildCmdShell: bc.RuntimeCmd.BuildCmdShell,
		}
	case *pb.BuildConfig_RuntimeDockerfile:
		return &domain.BuildConfigRuntimeDockerfile{
			DockerfileName: bc.RuntimeDockerfile.DockerfileName,
			Context:        bc.RuntimeDockerfile.Context,
		}
	case *pb.BuildConfig_StaticCmd:
		return &domain.BuildConfigStaticCmd{
			BaseImage:     bc.StaticCmd.BaseImage,
			BuildCmd:      bc.StaticCmd.BuildCmd,
			BuildCmdShell: bc.StaticCmd.BuildCmdShell,
			ArtifactPath:  bc.StaticCmd.ArtifactPath,
		}
	case *pb.BuildConfig_StaticDockerfile:
		return &domain.BuildConfigStaticDockerfile{
			DockerfileName: bc.StaticDockerfile.DockerfileName,
			Context:        bc.StaticDockerfile.Context,
			ArtifactPath:   bc.StaticDockerfile.ArtifactPath,
		}
	default:
		panic("unknown pb build config type")
	}
}

func ToPBBuildConfig(c domain.BuildConfig) *pb.BuildConfig {
	switch bc := c.(type) {
	case *domain.BuildConfigRuntimeCmd:
		return &pb.BuildConfig{BuildConfig: &pb.BuildConfig_RuntimeCmd{RuntimeCmd: &pb.BuildConfigRuntimeCmd{
			BaseImage:     bc.BaseImage,
			BuildCmd:      bc.BuildCmd,
			BuildCmdShell: bc.BuildCmdShell,
		}}}
	case *domain.BuildConfigRuntimeDockerfile:
		return &pb.BuildConfig{BuildConfig: &pb.BuildConfig_RuntimeDockerfile{RuntimeDockerfile: &pb.BuildConfigRuntimeDockerfile{
			DockerfileName: bc.DockerfileName,
			Context:        bc.Context,
		}}}
	case *domain.BuildConfigStaticCmd:
		return &pb.BuildConfig{BuildConfig: &pb.BuildConfig_StaticCmd{StaticCmd: &pb.BuildConfigStaticCmd{
			BaseImage:     bc.BaseImage,
			BuildCmd:      bc.BuildCmd,
			BuildCmdShell: bc.BuildCmdShell,
			ArtifactPath:  bc.ArtifactPath,
		}}}
	case *domain.BuildConfigStaticDockerfile:
		return &pb.BuildConfig{BuildConfig: &pb.BuildConfig_StaticDockerfile{StaticDockerfile: &pb.BuildConfigStaticDockerfile{
			DockerfileName: bc.DockerfileName,
			Context:        bc.Context,
			ArtifactPath:   bc.ArtifactPath,
		}}}
	default:
		panic("unknown domain build config type")
	}
}
