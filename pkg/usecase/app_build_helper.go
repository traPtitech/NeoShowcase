package usecase

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
)

type AppBuildHelper struct {
	component domain.ComponentService
	image     builder.ImageConfig
}

func NewAppBuildHelper(
	component domain.ComponentService,
	imageConfig builder.ImageConfig,
) *AppBuildHelper {
	return &AppBuildHelper{
		component: component,
		image:     imageConfig,
	}
}

func (s *AppBuildHelper) tryStartBuild(app *domain.Application, build *domain.Build) {
	switch app.DeployType {
	case domain.DeployTypeRuntime:
		s.component.BroadcastBuilder(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD_IMAGE,
			Body: &pb.BuilderRequest_BuildImage{
				BuildImage: &pb.StartBuildImageRequest{
					ImageName: s.image.FullImageName(app.ID),
					ImageTag:  build.Commit,
					Source: &pb.BuildSource{
						RepositoryId: app.RepositoryID,
						Commit:       build.Commit,
					},
					Options: &pb.BuildOptions{
						BaseImageName:  app.Config.BaseImage,
						DockerfileName: app.Config.DockerfileName,
						ArtifactPath:   app.Config.ArtifactPath,
						BuildCmd:       app.Config.BuildCmd,
						EntrypointCmd:  app.Config.EntrypointCmd,
					},
					BuildId:       build.ID,
					ApplicationId: app.ID,
				},
			},
		})

	case domain.DeployTypeStatic:
		s.component.BroadcastBuilder(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD_STATIC,
			Body: &pb.BuilderRequest_BuildStatic{
				BuildStatic: &pb.StartBuildStaticRequest{
					Source: &pb.BuildSource{
						RepositoryId: app.RepositoryID,
						Commit:       build.Commit,
					},
					Options: &pb.BuildOptions{
						BaseImageName:  app.Config.BaseImage,
						DockerfileName: app.Config.DockerfileName,
						ArtifactPath:   app.Config.ArtifactPath,
						BuildCmd:       app.Config.BuildCmd,
						EntrypointCmd:  app.Config.EntrypointCmd,
					},
					BuildId:       build.ID,
					ApplicationId: app.ID,
				},
			},
		})
	}
}
