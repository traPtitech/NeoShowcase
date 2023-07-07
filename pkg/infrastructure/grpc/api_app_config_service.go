package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetEnvVars(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.ApplicationEnvVars], error) {
	environments, err := s.svc.GetEnvironmentVariables(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ApplicationEnvVars{
		Variables: ds.Map(environments, pbconvert.ToPBEnvironment),
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

func (s *APIService) DeleteEnvVar(ctx context.Context, req *connect.Request[pb.DeleteApplicationEnvVarRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	err := s.svc.DeleteEnvironmentVariable(ctx, msg.ApplicationId, msg.Key)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) GetOutput(ctx context.Context, req *connect.Request[pb.GetOutputRequest]) (*connect.Response[pb.ApplicationOutputs], error) {
	msg := req.Msg
	if msg.Before == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("before cannot be null"))
	}
	before := msg.Before.AsTime()
	logs, err := s.svc.GetOutput(ctx, msg.ApplicationId, before)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ApplicationOutputs{
		Outputs: ds.Map(logs, pbconvert.ToPBApplicationOutput),
	})
	return res, nil
}

func (s *APIService) GetOutputStream(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest], st *connect.ServerStream[pb.ApplicationOutput]) error {
	err := s.svc.GetOutputStream(ctx, req.Msg.Id, func(l *domain.ContainerLog) error {
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
