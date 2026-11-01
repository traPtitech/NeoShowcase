package grpc

import (
	"context"
	"time"

	"connectrpc.com/connect/v2"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetAvailableMetrics(ctx context.Context, _ *emptypb.Empty) (*pb.AvailableMetrics, error) {
	names := s.svc.GetAvailableMetrics(ctx)
	res := &pb.AvailableMetrics{
		MetricsNames: names,
	}
	return res, nil
}

func (s *APIService) GetApplicationMetrics(ctx context.Context, req *pb.GetApplicationMetricsRequest) (*pb.ApplicationMetrics, error) {
	if req.Before == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, "before cannot be null")
	}
	metrics, err := s.svc.GetApplicationMetrics(ctx,
		req.MetricsName,
		req.ApplicationId,
		req.Before.AsTime(),
		time.Duration(req.LimitSeconds)*time.Second,
	)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.ApplicationMetrics{
		Metrics: ds.Map(metrics, pbconvert.ToPBApplicationMetric),
	}
	return res, nil
}

func (s *APIService) GetOutput(ctx context.Context, req *pb.GetOutputRequest) (*pb.ApplicationOutputs, error) {
	if req.Before == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, "before cannot be null")
	}
	before := req.Before.AsTime()
	logs, err := s.svc.GetOutput(ctx, req.ApplicationId, before, int(req.Limit))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.ApplicationOutputs{
		Outputs: ds.Map(logs, pbconvert.ToPBApplicationOutput),
	}
	return res, nil
}

func (s *APIService) GetOutputStream(ctx context.Context, req *pb.GetOutputStreamRequest, st pbconnect.APIServiceGetOutputStreamServerStream) error {
	if req.Begin == nil {
		return connect.NewError(connect.CodeInvalidArgument, "begin cannot be null")
	}
	begin := req.Begin.AsTime()
	err := s.svc.GetOutputStream(ctx, req.ApplicationId, begin, func(l *domain.ContainerLog) error {
		return st.Send(pbconvert.ToPBApplicationOutput(l))
	})
	if err != nil {
		return handleUseCaseError(err)
	}
	return nil
}
