package grpc

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
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
		DeployType:       config.BuildConfig.BuildType().DeployType(),
		Running:          msg.StartOnCreate,
		Container:        domain.ContainerStateMissing,
		CurrentCommit:    domain.EmptyCommit,
		WantCommit:       domain.EmptyCommit,
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
	res := connect.NewResponse(pbconvert.ToPBApplication(app))
	return res, nil
}

func (s *APIService) GetApplications(ctx context.Context, req *connect.Request[pb.GetApplicationsRequest]) (*connect.Response[pb.GetApplicationsResponse], error) {
	all := req.Msg.Scope == pb.GetApplicationsRequest_ALL
	applications, err := s.svc.GetApplications(ctx, all)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetApplicationsResponse{
		Applications: ds.Map(applications, pbconvert.ToPBApplication),
	})
	return res, nil
}

func (s *APIService) GetApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.Application], error) {
	application, err := s.svc.GetApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBApplication(application))
	return res, nil
}

func (s *APIService) UpdateApplication(ctx context.Context, req *connect.Request[pb.UpdateApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	app, err := s.svc.GetApplication(ctx, msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}

	websites := app.Websites
	for _, createReq := range msg.NewWebsites {
		websites = append(websites, pbconvert.FromPBCreateWebsiteRequest(createReq))
	}
	for _, deleteReq := range msg.DeleteWebsites {
		websites = lo.Reject(websites, func(w *domain.Website, _ int) bool { return w.ID == deleteReq.Id })
	}

	err = s.svc.UpdateApplication(ctx, msg.Id, &domain.UpdateApplicationArgs{
		Name:             optional.FromNonZero(msg.Name),
		RefName:          optional.FromNonZero(msg.RefName),
		UpdatedAt:        optional.From(time.Now()),
		Config:           optional.Map(optional.FromNonZero(msg.Config), pbconvert.FromPBApplicationConfig),
		Websites:         optional.FromNonZeroSlice(websites),
		PortPublications: optional.FromNonZeroSlice(ds.Map(msg.PortPublications, pbconvert.FromPBPortPublication)),
		OwnerIDs:         optional.FromNonZeroSlice(msg.OwnerIds),
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
