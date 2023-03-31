package usecase

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
)

type AppBuildService interface {
	TryStartBuild(app *domain.Application, build *domain.Build)
}

type appBuildService struct {
	buildRepo domain.BuildRepository
	component domain.ComponentService
	image     builder.ImageConfig
}

func NewAppBuildService(
	buildRepo domain.BuildRepository,
	component domain.ComponentService,
	imageConfig builder.ImageConfig,
) AppBuildService {
	return &appBuildService{
		buildRepo: buildRepo,
		component: component,
		image:     imageConfig,
	}
}

func (s *appBuildService) TryStartBuild(app *domain.Application, build *domain.Build) {
	switch app.BuildType {
	case builder.BuildTypeRuntime:
		s.component.BroadcastBuilder(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD_IMAGE,
			Body: &pb.BuilderRequest_BuildImage{
				BuildImage: &pb.StartBuildImageRequest{
					ImageName: s.image.ImageName(app.ID),
					ImageTag:  build.ID,
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

	case builder.BuildTypeStatic:
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
