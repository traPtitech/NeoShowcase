package grpc

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/friendsofgo/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetAvailableMetrics(ctx context.Context, c *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailableMetrics], error) {
	names := s.svc.GetAvailableMetrics(ctx)
	res := connect.NewResponse(&pb.AvailableMetrics{
		MetricsNames: names,
	})
	return res, nil
}

func (s *APIService) GetApplicationMetrics(ctx context.Context, req *connect.Request[pb.GetApplicationMetricsRequest]) (*connect.Response[pb.ApplicationMetrics], error) {
	msg := req.Msg
	if msg.Before == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("before cannot be null"))
	}
	metrics, err := s.svc.GetApplicationMetrics(ctx,
		msg.MetricsName,
		msg.ApplicationId,
		msg.Before.AsTime(),
		time.Duration(msg.LimitSeconds)*time.Second,
	)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ApplicationMetrics{
		Metrics: ds.Map(metrics, pbconvert.ToPBApplicationMetric),
	})
	return res, nil
}

func (s *APIService) GetOutput(ctx context.Context, req *connect.Request[pb.GetOutputRequest]) (*connect.Response[pb.ApplicationOutputs], error) {
	msg := req.Msg
	if msg.Before == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("before cannot be null"))
	}
	before := msg.Before.AsTime()
	logs, err := s.svc.GetOutput(ctx, msg.ApplicationId, before, int(msg.Limit))
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
