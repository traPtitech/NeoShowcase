package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var AppScopeMapper = mapper.MustNewValueMapper(map[pb.GetApplicationsRequest_Scope]apiserver.GetAppScopeType{
	pb.GetApplicationsRequest_MINE:       apiserver.GetAppScopeMine,
	pb.GetApplicationsRequest_ALL:        apiserver.GetAppScopeAll,
	pb.GetApplicationsRequest_REPOSITORY: apiserver.GetAppScopeRepository,
})

var DeployTypeMapper = mapper.MustNewValueMapper(map[domain.DeployType]pb.DeployType{
	domain.DeployTypeRuntime: pb.DeployType_RUNTIME,
	domain.DeployTypeStatic:  pb.DeployType_STATIC,
})

var ContainerStateMapper = mapper.MustNewValueMapper(map[domain.ContainerState]pb.Application_ContainerState{
	domain.ContainerStateMissing:    pb.Application_MISSING,
	domain.ContainerStateStarting:   pb.Application_STARTING,
	domain.ContainerStateRestarting: pb.Application_RESTARTING,
	domain.ContainerStateRunning:    pb.Application_RUNNING,
	domain.ContainerStateExited:     pb.Application_EXITED,
	domain.ContainerStateErrored:    pb.Application_ERRORED,
	domain.ContainerStateUnknown:    pb.Application_UNKNOWN,
})

func ToPBApplication(app *domain.Application, latestBuild *domain.Build) *pb.Application {
	pbApp := &pb.Application{
		Id:               app.ID,
		Name:             app.Name,
		RepositoryId:     app.RepositoryID,
		RefName:          app.RefName,
		Commit:           app.Commit,
		DeployType:       DeployTypeMapper.IntoMust(app.DeployType),
		Running:          app.Running,
		Container:        ContainerStateMapper.IntoMust(app.Container),
		ContainerMessage: app.ContainerMessage,
		CurrentBuild:     app.CurrentBuild,
		CreatedAt:        timestamppb.New(app.CreatedAt),
		UpdatedAt:        timestamppb.New(app.UpdatedAt),
		Config:           ToPBApplicationConfig(app.Config),
		Websites:         ds.Map(app.Websites, ToPBWebsite),
		PortPublications: ds.Map(app.PortPublications, ToPBPortPublication),
		OwnerIds:         app.OwnerIDs,
	}
	if latestBuild != nil {
		status := BuildStatusMapper.IntoMust(latestBuild.Status)
		pbApp.LatestBuildStatus = &status
	}
	return pbApp
}

func FromPBApplication(app *pb.Application) *domain.Application {
	return &domain.Application{
		ID:               app.Id,
		Name:             app.Name,
		RepositoryID:     app.RepositoryId,
		RefName:          app.RefName,
		Commit:           app.Commit,
		DeployType:       DeployTypeMapper.FromMust(app.DeployType),
		Running:          app.Running,
		Container:        ContainerStateMapper.FromMust(app.Container),
		ContainerMessage: app.ContainerMessage,
		CurrentBuild:     app.CurrentBuild,
		CreatedAt:        app.CreatedAt.AsTime(),
		UpdatedAt:        app.UpdatedAt.AsTime(),
		Config:           FromPBApplicationConfig(app.Config),
		Websites:         ds.Map(app.Websites, FromPBWebsite),
		PortPublications: ds.Map(app.PortPublications, FromPBPortPublication),
		OwnerIDs:         app.OwnerIds,
	}
}

func FromPBUpdateOwners(req *pb.UpdateApplicationRequest_UpdateOwners) []string {
	return req.OwnerIds
}
