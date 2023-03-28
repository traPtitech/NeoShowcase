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
	buildRepo       domain.BuildRepository
	component       domain.ComponentService
	imageRegistry   string
	imageNamePrefix string
}

func NewAppBuildService(
	buildRepo domain.BuildRepository,
	component domain.ComponentService,
	registry builder.DockerImageRegistryString,
	prefix builder.DockerImageNamePrefixString,
) AppBuildService {
	return &appBuildService{
		buildRepo:       buildRepo,
		component:       component,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
}

func (s *appBuildService) TryStartBuild(app *domain.Application, build *domain.Build) {
	switch app.BuildType {
	case builder.BuildTypeRuntime:
		s.component.TryStartBuild(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD_IMAGE,
			Body: &pb.BuilderRequest_BuildImage{
				BuildImage: &pb.StartBuildImageRequest{
					ImageName: builder.GetImageName(s.imageRegistry, s.imageNamePrefix, app.ID),
					ImageTag:  build.ID,
					Source: &pb.BuildSource{
						RepositoryUrl: app.Repository.URL,
						Commit:        build.Commit,
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
		s.component.TryStartBuild(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD_STATIC,
			Body: &pb.BuilderRequest_BuildStatic{
				BuildStatic: &pb.StartBuildStaticRequest{
					Source: &pb.BuildSource{
						RepositoryUrl: app.Repository.URL,
						Commit:        build.Commit,
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
