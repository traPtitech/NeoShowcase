package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var DeployTypeMapper = mapper.MustNewValueMapper(map[domain.DeployType]pb.DeployType{
	domain.DeployTypeRuntime: pb.DeployType_RUNTIME,
	domain.DeployTypeStatic:  pb.DeployType_STATIC,
})

var ContainerStateMapper = mapper.MustNewValueMapper(map[domain.ContainerState]pb.Application_ContainerState{
	domain.ContainerStateMissing:  pb.Application_MISSING,
	domain.ContainerStateStarting: pb.Application_STARTING,
	domain.ContainerStateRunning:  pb.Application_RUNNING,
	domain.ContainerStateExited:   pb.Application_EXITED,
	domain.ContainerStateErrored:  pb.Application_ERRORED,
	domain.ContainerStateUnknown:  pb.Application_UNKNOWN,
})

func ToPBApplication(app *domain.Application) *pb.Application {
	return &pb.Application{
		Id:               app.ID,
		Name:             app.Name,
		RepositoryId:     app.RepositoryID,
		RefName:          app.RefName,
		DeployType:       DeployTypeMapper.IntoMust(app.DeployType),
		Running:          app.Running,
		Container:        ContainerStateMapper.IntoMust(app.Container),
		CurrentCommit:    app.CurrentCommit,
		WantCommit:       app.WantCommit,
		CreatedAt:        timestamppb.New(app.CreatedAt),
		UpdatedAt:        timestamppb.New(app.UpdatedAt),
		Config:           ToPBApplicationConfig(app.Config),
		Websites:         ds.Map(app.Websites, ToPBWebsite),
		PortPublications: ds.Map(app.PortPublications, ToPBPortPublication),
		OwnerIds:         app.OwnerIDs,
	}
}

func FromPBUpdateOwners(req *pb.UpdateApplicationRequest_UpdateOwners) []string {
	return req.OwnerIds
}
