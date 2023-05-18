package usecase

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

type AppBuildHelper struct {
	builder domain.ControllerBuilderService
	image   builder.ImageConfig
}

func NewAppBuildHelper(
	builder domain.ControllerBuilderService,
	imageConfig builder.ImageConfig,
) *AppBuildHelper {
	return &AppBuildHelper{
		builder: builder,
		image:   imageConfig,
	}
}

func (s *AppBuildHelper) tryStartBuild(app *domain.Application, build *domain.Build) {
	s.builder.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_START_BUILD,
		Body: &pb.BuilderRequest_StartBuild{StartBuild: &pb.StartBuildRequest{
			ApplicationId: app.ID,
			BuildId:       build.ID,
			RepositoryId:  app.RepositoryID,
			Commit:        build.Commit,
			ImageName:     s.image.FullImageName(app.ID),
			ImageTag:      build.Commit,
			Config:        pbconvert.ToPBApplicationConfig(app.Config),
		}},
	})
}
