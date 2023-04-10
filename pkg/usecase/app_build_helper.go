package usecase

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pbconvert"
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
	s.component.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_START_BUILD,
		Body: &pb.BuilderRequest_StartBuild{StartBuild: &pb.StartBuildRequest{
			ApplicationId: app.ID,
			BuildId:       build.ID,
			RepositoryId:  app.RepositoryID,
			Commit:        build.Commit,
			ImageName:     s.image.FullImageName(app.ID),
			ImageTag:      build.Commit,
			BuildConfig:   pbconvert.ToPBBuildConfig(app.Config.BuildConfig),
		}},
	})
}
