package grpc

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StaticSiteService struct {
	svc usecase.StaticSiteService

	pb.UnimplementedStaticSiteServiceServer
}

func NewStaticSiteServiceServer(svc usecase.StaticSiteService) *StaticSiteService {
	return &StaticSiteService{svc: svc}
}

func (s *StaticSiteService) Reload(ctx context.Context, request *pb.ReloadRequest) (*pb.ReloadResponse, error) {
	err := s.svc.Reload(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ReloadResponse{}, nil
}
