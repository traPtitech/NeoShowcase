package grpc

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/samber/lo"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIService) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.Application, error) {
	user := web.GetUser(ctx)
	now := time.Now()
	config := pbconvert.FromPBApplicationConfig(req.Config)

	repo, err := s.svc.GetRepository(ctx, req.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	ownerIDs := lo.Uniq(append(repo.OwnerIDs, user.ID))

	app := &domain.Application{
		ID:               domain.NewID(),
		Name:             req.Name,
		RepositoryID:     req.RepositoryId,
		RefName:          req.RefName,
		Commit:           domain.EmptyCommit,
		DeployType:       config.BuildConfig.BuildType().DeployType(),
		Running:          req.StartOnCreate,
		Container:        domain.ContainerStateMissing,
		CurrentBuild:     "",
		CreatedAt:        now,
		UpdatedAt:        now,
		Config:           config,
		Websites:         ds.Map(req.Websites, pbconvert.FromPBCreateWebsiteRequest),
		PortPublications: ds.Map(req.PortPublications, pbconvert.FromPBPortPublication),
		OwnerIDs:         ownerIDs,
	}
	app, err = s.svc.CreateApplication(ctx, app)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return pbconvert.ToPBApplication(app, nil), nil
}

func (s *APIService) GetApplications(ctx context.Context, req *pb.GetApplicationsRequest) (*pb.GetApplicationsResponse, error) {
	scope := apiserver.GetAppScope{
		Scope:        pbconvert.AppScopeMapper.IntoMust(req.Scope),
		RepositoryID: optional.FromPtr(req.RepositoryId),
	}
	apps, err := s.svc.GetApplications(ctx, scope)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &pb.GetApplicationsResponse{
		Applications: ds.Map(apps, func(app *apiserver.TopAppInfo) *pb.Application {
			return pbconvert.ToPBApplication(app.App, app.LatestBuild)
		}),
	}, nil
}

func (s *APIService) GetApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.Application, error) {
	app, err := s.svc.GetApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return pbconvert.ToPBApplication(app.App, app.LatestBuild), nil
}

func (s *APIService) UpdateApplication(ctx context.Context, req *pb.UpdateApplicationRequest) (*emptypb.Empty, error) {
	err := s.svc.UpdateApplication(ctx, req.Id, &domain.UpdateApplicationArgs{
		Name:             optional.FromPtr(req.Name),
		RefName:          optional.FromPtr(req.RefName),
		UpdatedAt:        optional.From(time.Now()),
		Config:           optional.FromNonZero(req.Config).Map(pbconvert.FromPBApplicationConfig),
		Websites:         optional.FromNonZero(req.Websites).Map(pbconvert.FromPBUpdateWebsites),
		PortPublications: optional.FromNonZero(req.PortPublications).Map(pbconvert.FromPBUpdatePorts),
		OwnerIDs:         optional.FromNonZero(req.OwnerIds).Map(pbconvert.FromPBUpdateOwners),
	})
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *APIService) DeleteApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}
