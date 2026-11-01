package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetEnvVars(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.ApplicationEnvVars, error) {
	environments, err := s.svc.GetEnvironmentVariables(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.ApplicationEnvVars{
		Variables: ds.Map(environments, pbconvert.ToPBEnvironment),
	}
	return res, nil
}

func (s *APIService) SetEnvVar(ctx context.Context, req *pb.SetApplicationEnvVarRequest) (*emptypb.Empty, error) {
	msg := req
	err := s.svc.SetEnvironmentVariable(ctx, msg.ApplicationId, msg.Key, msg.Value)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) DeleteEnvVar(ctx context.Context, req *pb.DeleteApplicationEnvVarRequest) (*emptypb.Empty, error) {
	msg := req
	err := s.svc.DeleteEnvironmentVariable(ctx, msg.ApplicationId, msg.Key)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) StartApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.StartApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) StopApplication(ctx context.Context, req *pb.ApplicationIdRequest) (*emptypb.Empty, error) {
	err := s.svc.StopApplication(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}
