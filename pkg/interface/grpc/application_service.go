package grpc

import (
	"context"
	"fmt"

	"github.com/friendsofgo/errors"
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

func handleUseCaseError(err error) error {
	var uErr *usecase.Error
	if errors.As(err, &uErr) {
		switch uErr.Type {
		case usecase.ErrorTypeBadRequest:
			return status.Errorf(codes.InvalidArgument, "%v", err)
		case usecase.ErrorTypeNotFound:
			return status.Errorf(codes.NotFound, "%v", err)
		case usecase.ErrorTypeAlreadyExists:
			return status.Errorf(codes.AlreadyExists, "%v", err)
		}
	}
	return status.Errorf(codes.Internal, "%v", err)
}

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
		return nil, handleUseCaseError(err)
	}
	return &pb.GetApplicationsResponse{
		Applications: lo.Map(applications, func(app *domain.Application, i int) *pb.Application {
			return toPBApplication(app)
		}),
	}, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.Application, error) {
	application, err := s.svc.CreateApplication(ctx, usecase.CreateApplicationArgs{
		UserID:        getUserID(),
		Name:          req.Name,
		RepositoryURL: req.RepositoryUrl,
		BranchName:    req.BranchName,
		BuildType:     fromPBBuildType(req.BuildType),
		Config:        fromPBApplicationConfig(req.Config),
		Websites: lo.Map(req.Websites, func(website *pb.CreateWebsiteRequest, i int) *domain.Website {
			return &domain.Website{
				ID:         domain.NewID(),
				FQDN:       website.Fqdn,
				PathPrefix: website.PathPrefix,
				HTTPS:      website.Https,
				HTTPPort:   int(website.HttpPort),
			}
		}),
		StartOnCreate: req.StartOnCreate,
	})
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return toPBApplication(application), nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.Application, error) {
	application, err := s.svc.GetApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return toPBApplication(application), nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) GetApplicationBuilds(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.GetApplicationBuildsResponse, error) {
	builds, err := s.svc.GetApplicationBuilds(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &pb.GetApplicationBuildsResponse{
		Builds: lo.Map(builds, func(build *domain.Build, i int) *pb.Build {
			return toPBBuild(build)
		}),
	}, nil
}

func (s *ApplicationService) GetApplicationBuild(ctx context.Context, req *pb.GetApplicationBuildRequest) (*pb.Build, error) {
	build, err := s.svc.GetApplicationBuild(ctx, req.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return toPBBuild(build), nil
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
		return nil, handleUseCaseError(err)
	}
	return &pb.ApplicationEnvironmentVariables{
		Variables: lo.Map(environments, func(env *domain.Environment, i int) *pb.ApplicationEnvironmentVariable {
			return toPBEnvironment(env)
		}),
	}, nil
}

func (s *ApplicationService) SetApplicationEnvironmentVariable(ctx context.Context, req *pb.SetApplicationEnvironmentVariableRequest) (*emptypb.Empty, error) {
	err := s.svc.SetApplicationEnvironmentVariable(ctx, req.ApplicationId, req.Key, req.Value)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) GetApplicationOutput(context.Context, *pb.ApplicationIdRequest) (*pb.ApplicationOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationOutput not implemented")
}

func (s *ApplicationService) GetApplicationKeys(context.Context, *pb.ApplicationIdRequest) (*pb.ApplicationKeys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationKeys not implemented")
}

func (s *ApplicationService) RetryCommitBuild(ctx context.Context, req *pb.RetryCommitBuildRequest) (*emptypb.Empty, error) {
	err := s.svc.RetryCommitBuild(ctx, req.ApplicationId, req.Commit)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) StartApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.StartApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ApplicationService) StopApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.StopApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func fromPBBuildType(buildType pb.BuildType) builder.BuildType {
	switch buildType {
	case pb.BuildType_RUNTIME:
		return builder.BuildTypeRuntime
	case pb.BuildType_STATIC:
		return builder.BuildTypeStatic
	default:
		panic(fmt.Sprintf("unknown build type: %v", buildType))
	}
}

func toPBBuildType(buildType builder.BuildType) pb.BuildType {
	switch buildType {
	case builder.BuildTypeRuntime:
		return pb.BuildType_RUNTIME
	case builder.BuildTypeStatic:
		return pb.BuildType_STATIC
	default:
		panic(fmt.Sprintf("unknown build type: %v", buildType))
	}
}

func toPBApplicationState(state domain.ApplicationState) pb.ApplicationState {
	switch state {
	case domain.ApplicationStateIdle:
		return pb.ApplicationState_IDLE
	case domain.ApplicationStateDeploying:
		return pb.ApplicationState_DEPLOYING
	case domain.ApplicationStateRunning:
		return pb.ApplicationState_RUNNING
	case domain.ApplicationStateErrored:
		return pb.ApplicationState_ERRORED
	default:
		panic(fmt.Sprintf("unknown application state: %v", state))
	}
}

func toPBAuthenticationType(t domain.AuthenticationType) pb.AuthenticationType {
	switch t {
	case domain.AuthenticationTypeOff:
		return pb.AuthenticationType_OFF
	case domain.AuthenticationTypeSoft:
		return pb.AuthenticationType_SOFT
	case domain.AuthenticationTypeHard:
		return pb.AuthenticationType_HARD
	default:
		panic(fmt.Sprintf("unknown authentication type: %v", t))
	}
}

func fromPBAuthenticationType(t pb.AuthenticationType) domain.AuthenticationType {
	switch t {
	case pb.AuthenticationType_OFF:
		return domain.AuthenticationTypeOff
	case pb.AuthenticationType_SOFT:
		return domain.AuthenticationTypeSoft
	case pb.AuthenticationType_HARD:
		return domain.AuthenticationTypeHard
	default:
		panic(fmt.Sprintf("unknown authentication type: %v", t))
	}
}

func toPBApplicationConfig(c domain.ApplicationConfig) *pb.ApplicationConfig {
	return &pb.ApplicationConfig{
		UseMariadb:     c.UseMariaDB,
		UseMongodb:     c.UseMongoDB,
		BaseImage:      c.BaseImage,
		DockerfileName: c.DockerfileName,
		ArtifactPath:   c.ArtifactPath,
		BuildCmd:       c.BuildCmd,
		EntrypointCmd:  c.EntrypointCmd,
		Authentication: toPBAuthenticationType(c.Authentication),
	}
}

func fromPBApplicationConfig(c *pb.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		UseMariaDB:     c.UseMariadb,
		UseMongoDB:     c.UseMongodb,
		BaseImage:      c.BaseImage,
		DockerfileName: c.DockerfileName,
		ArtifactPath:   c.ArtifactPath,
		BuildCmd:       c.BuildCmd,
		EntrypointCmd:  c.EntrypointCmd,
		Authentication: fromPBAuthenticationType(c.Authentication),
	}
}

func toPBWebsite(website *domain.Website) *pb.Website {
	return &pb.Website{
		Id:         website.ID,
		Fqdn:       website.FQDN,
		PathPrefix: website.PathPrefix,
		Https:      website.HTTPS,
		HttpPort:   int32(website.HTTPPort),
	}
}

func toPBApplication(app *domain.Application) *pb.Application {
	return &pb.Application{
		Id:            app.ID,
		Name:          app.Name,
		RepositoryUrl: app.Repository.URL,
		BranchName:    app.BranchName,
		BuildType:     toPBBuildType(app.BuildType),
		State:         toPBApplicationState(app.State),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
		Config:        toPBApplicationConfig(app.Config),
		Websites:      lo.Map(app.Websites, func(website *domain.Website, i int) *pb.Website { return toPBWebsite(website) }),
	}
}

func toPBBuildStatus(status builder.BuildStatus) pb.Build_BuildStatus {
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

func toPBBuild(build *domain.Build) *pb.Build {
	return &pb.Build{
		Id:        build.ID,
		Commit:    build.Commit,
		Status:    toPBBuildStatus(build.Status),
		StartedAt: timestamppb.New(build.StartedAt),
		FinishedAt: &pb.NullTimestamp{
			Timestamp: timestamppb.New(build.FinishedAt.V),
			Valid:     build.FinishedAt.Valid,
		},
		Retriable: build.Retriable,
	}
}

func toPBEnvironment(env *domain.Environment) *pb.ApplicationEnvironmentVariable {
	return &pb.ApplicationEnvironmentVariable{
		Key:   env.Key,
		Value: env.Value,
	}
}
