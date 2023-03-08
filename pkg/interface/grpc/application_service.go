package grpc

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type ApplicationService struct {
	svc usecase.APIServerService

	pb.UnimplementedApplicationServiceServer
}

func NewApplicationServiceServer(svc usecase.APIServerService) *ApplicationService {
	return &ApplicationService{
		svc: svc,
	}
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
	if err != nil {
		switch err {
		case usecase.ErrAlreadyExists:
			return nil, status.Errorf(codes.AlreadyExists, "app already exists")
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	return convertToPBApplication(application), nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.Application, error) {
	application, err := s.svc.GetApplication(ctx, req.Id)
	if err != nil {
		if err == usecase.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "not found")
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return convertToPBApplication(application), nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteApplication(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) GetApplicationBuilds(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.GetApplicationBuildsResponse, error) {
	builds, err := s.svc.GetApplicationBuilds(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.GetApplicationBuildsResponse{
		Builds: lo.Map(builds, func(build *domain.Build, i int) *pb.Build {
			return convertToPBBuild(build)
		}),
	}, nil
}

func (s *ApplicationService) GetApplicationBuild(ctx context.Context, req *pb.GetApplicationBuildRequest) (*pb.Build, error) {
	build, err := s.svc.GetApplicationBuild(ctx, req.BuildId)
	if err != nil {
		if err == usecase.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "not found")
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return convertToPBBuild(build), nil
}

func (s *ApplicationService) GetApplicationBuildLog(context.Context, *pb.GetApplicationBuildLogRequest) (*pb.BuildLog, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationBuildLog not implemented")
}

func (s *ApplicationService) GetApplicationBuildArtifact(context.Context, *pb.ApplicationIdRequest) (*pb.ApplicationBuildArtifact, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationBuildArtifact not implemented")
}

func (s *ApplicationService) GetApplicationEnvironmentVariables(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.ApplicationEnvironmentVariables, error) {
	environments, err := s.svc.GetApplicationEnvironmentVariables(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.ApplicationEnvironmentVariables{
		Variables: lo.Map(environments, func(env *domain.Environment, i int) *pb.ApplicationEnvironmentVariable {
			return convertToPBEnvironment(env)
		}),
	}, nil
}

func (s *ApplicationService) SetApplicationEnvironmentVariable(ctx context.Context, req *pb.SetApplicationEnvironmentVariableRequest) (*emptypb.Empty, error) {
	err := s.svc.SetApplicationEnvironmentVariable(ctx, req.ApplicationId, req.Key, req.Value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) GetApplicationOutput(context.Context, *pb.ApplicationIdRequest) (*pb.ApplicationOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationOutput not implemented")
}

func (s *ApplicationService) GetApplicationKeys(context.Context, *pb.ApplicationIdRequest) (*pb.ApplicationKeys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationKeys not implemented")
}

func (s *ApplicationService) StartApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.StartApplication(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) RestartApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.RestartApplication(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) StopApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.StopApplication(ctx, req.Id)
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

func convertToPBBuildStatus(status builder.BuildStatus) pb.Build_BuildStatus {
	switch status {
	case builder.BuildStatusBuilding:
		return pb.Build_BUILDING
	case builder.BuildStatusSucceeded:
		return pb.Build_SUCCEEDED
	case builder.BuildStatusFailed:
		return pb.Build_FAILED
	case builder.BuildStatusCanceled:
		return pb.Build_CANCELLED
	case builder.BuildStatusQueued:
		return pb.Build_QUEUED
	case builder.BuildStatusSkipped:
		return pb.Build_SKIPPED
	default:
		panic(fmt.Sprintf("unknown build status: %v", status))
	}
}

func convertToPBBuild(build *domain.Build) *pb.Build {
	return &pb.Build{
		Id:        build.ID,
		Status:    convertToPBBuildStatus(build.Status),
		StartedAt: timestamppb.New(build.StartedAt),
		FinishedAt: &pb.NullTimestamp{
			Timestamp: timestamppb.New(build.FinishedAt.V),
			Valid:     build.FinishedAt.Valid,
		},
	}
}

func convertToPBEnvironment(env *domain.Environment) *pb.ApplicationEnvironmentVariable {
	return &pb.ApplicationEnvironmentVariable{
		Key:   env.Key,
		Value: env.Value,
	}
}
