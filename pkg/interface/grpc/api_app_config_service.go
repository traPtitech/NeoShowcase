package grpc

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pbconvert"
)

func (s *APIService) GetEnvVars(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.ApplicationEnvVars], error) {
	environments, err := s.svc.GetEnvironmentVariables(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ApplicationEnvVars{
		Variables: lo.Map(environments, func(env *domain.Environment, i int) *pb.ApplicationEnvVar {
			return pbconvert.ToPBEnvironment(env)
		}),
	})
	return res, nil
}

func (s *APIService) SetEnvVar(ctx context.Context, req *connect.Request[pb.SetApplicationEnvVarRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	err := s.svc.SetEnvironmentVariable(ctx, msg.ApplicationId, msg.Key, msg.Value)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) GetOutput(ctx context.Context, req *connect.Request[pb.GetOutputRequest]) (*connect.Response[pb.GetOutputResponse], error) {
	msg := req.Msg
	before := time.Now()
	if req.Msg.Before != nil {
		before = msg.Before.AsTime()
	}
	logs, err := s.svc.GetOutput(ctx, msg.ApplicationId, before)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetOutputResponse{
		Outputs: lo.Map(logs, func(l *domain.ContainerLog, i int) *pb.ApplicationOutput {
			return pbconvert.ToPBApplicationOutput(l)
		}),
	})
	return res, nil
}

func (s *APIService) GetOutputStream(ctx context.Context, req *connect.Request[pb.GetOutputStreamRequest], st *connect.ServerStream[pb.ApplicationOutput]) error {
	msg := req.Msg
	after := time.Now()
	if req.Msg.After != nil {
		after = msg.After.AsTime()
	}
	err := s.svc.GetOutputStream(ctx, msg.ApplicationId, after, func(l *domain.ContainerLog) error {
		return st.Send(pbconvert.ToPBApplicationOutput(l))
	})
	if err != nil {
		return handleUseCaseError(err)
	}
	return nil
}

func (s *APIService) StartApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.StartApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) StopApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.StopApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
