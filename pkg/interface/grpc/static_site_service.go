package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type StaticSiteService struct {
	svc usecase.StaticSiteServerService

	pb.UnimplementedStaticSiteServiceServer
}

func NewStaticSiteServiceServer(svc usecase.StaticSiteServerService) *StaticSiteService {
	return &StaticSiteService{svc: svc}
}

func (s *StaticSiteService) Reload(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.svc.Reload(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
