package grpc

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIService) CreateApplication(ctx context.Context, req *connect.Request[pb.CreateApplicationRequest]) (*connect.Response[pb.Application], error) {
	msg := req.Msg
	user := web.GetUser(ctx)
	now := time.Now()
	config := pbconvert.FromPBApplicationConfig(msg.Config)
	app := &domain.Application{
		ID:               domain.NewID(),
		Name:             msg.Name,
		RepositoryID:     msg.RepositoryId,
		RefName:          msg.RefName,
		Commit:           domain.EmptyCommit,
		DeployType:       config.BuildConfig.BuildType().DeployType(),
		Running:          msg.StartOnCreate,
		Container:        domain.ContainerStateMissing,
		CurrentBuild:     "",
		CreatedAt:        now,
		UpdatedAt:        now,
		Config:           config,
		Websites:         ds.Map(msg.Websites, pbconvert.FromPBCreateWebsiteRequest),
		PortPublications: ds.Map(msg.PortPublications, pbconvert.FromPBPortPublication),
		OwnerIDs:         []string{user.ID},
	}
	app, err := s.svc.CreateApplication(ctx, app)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBApplication(app, nil))
	return res, nil
}

func (s *APIService) GetApplications(ctx context.Context, req *connect.Request[pb.GetApplicationsRequest]) (*connect.Response[pb.GetApplicationsResponse], error) {
	all := req.Msg.Scope == pb.GetApplicationsRequest_ALL
	apps, err := s.svc.GetApplications(ctx, all)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetApplicationsResponse{
		Applications: ds.Map(apps, func(app *apiserver.TopAppInfo) *pb.Application {
			return pbconvert.ToPBApplication(app.App, app.LatestBuild)
		}),
	})
	return res, nil
}

func (s *APIService) GetApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.Application], error) {
	app, err := s.svc.GetApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBApplication(app.App, app.LatestBuild))
	return res, nil
}

func (s *APIService) UpdateApplication(ctx context.Context, req *connect.Request[pb.UpdateApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	err := s.svc.UpdateApplication(ctx, msg.Id, &domain.UpdateApplicationArgs{
		Name:             optional.FromPtr(msg.Name),
		RepositoryID:     optional.FromPtr(msg.RepositoryId),
		RefName:          optional.FromPtr(msg.RefName),
		UpdatedAt:        optional.From(time.Now()),
		Config:           optional.Map(optional.FromNonZero(msg.Config), pbconvert.FromPBApplicationConfig),
		Websites:         optional.Map(optional.FromNonZero(msg.Websites), pbconvert.FromPBUpdateWebsites),
		PortPublications: optional.Map(optional.FromNonZero(msg.PortPublications), pbconvert.FromPBUpdatePorts),
		OwnerIDs:         optional.Map(optional.FromNonZero(msg.OwnerIds), pbconvert.FromPBUpdateOwners),
	})
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) DeleteApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.DeleteApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
