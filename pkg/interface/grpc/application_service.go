package grpc

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type ApplicationService struct {
	svc usecase.APIServerService

	pb.UnimplementedApplicationServiceServer
}

func getUserID() string {
	return "tmp-user" // TODO: implement auth
}

func (s *ApplicationService) GetApplications(ctx context.Context, _ *emptypb.Empty) (*pb.GetApplicationsResponse, error) {
	applications, err := s.svc.GetApplicationsByUserID(ctx, getUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.GetApplicationsResponse{
		Applications: lo.Map(applications, func(app *domain.Application, i int) *pb.Application {
			return convertToPBApplication(app)
		}),
	}, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.Application, error) {
	application, err := s.svc.CreateApplication(ctx, usecase.CreateApplicationArgs{
		UserID:        getUserID(),
		RepositoryURL: req.RepositoryUrl,
		BranchName:    req.BranchName,
		BuildType:     convertFromPBBuildType(req.BuildType),
	})
	if err != nil { // TODO: handle possible user errors
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return convertToPBApplication(application), nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *pb.GetApplicationRequest) (*pb.Application, error) {
	application, err := s.svc.GetApplication(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return convertToPBApplication(application), nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, req *pb.DeleteApplicationRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteApplication(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

func convertFromPBBuildType(buildType pb.BuildType) builder.BuildType {
	switch buildType {
	case pb.BuildType_IMAGE:
		return builder.BuildTypeImage
	case pb.BuildType_STATIC:
		return builder.BuildTypeStatic
	default:
		panic(fmt.Sprintf("unknown build type: %v", buildType))
	}
}

func convertToPBBuildType(buildType builder.BuildType) pb.BuildType {
	switch buildType {
	case builder.BuildTypeImage:
		return pb.BuildType_IMAGE
	case builder.BuildTypeStatic:
		return pb.BuildType_STATIC
	default:
		panic(fmt.Sprintf("unknown build type: %v", buildType))
	}
}

func convertToPBApplication(app *domain.Application) *pb.Application {
	return &pb.Application{
		Id:            app.ID,
		RepositoryUrl: app.Repository.URL,
		BranchName:    app.BranchName,
		BuildType:     convertToPBBuildType(app.BuildType),
	}
}
