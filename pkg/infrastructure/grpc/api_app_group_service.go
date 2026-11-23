package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) CreateApplicationGroup(ctx context.Context, req *pb.CreateApplicationGroupRequest) (*pb.ApplicationGroup, error) {
	group, err := s.svc.CreateApplicationGroup(ctx, req.Name)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return pbconvert.ToPBApplicationGroup(group), nil
}

func (s *APIService) GetApplicationGroups(ctx context.Context, req *pb.GetApplicationGroupsRequest) (*pb.GetApplicationGroupsResponse, error) {
	groups, err := s.svc.GetApplicationGroups(ctx, pbconvert.AppGroupScopeMapper.FromMust(req.Scope))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &pb.GetApplicationGroupsResponse{
		Groups: ds.Map(groups, pbconvert.ToPBApplicationGroup),
	}, nil
}

func (s *APIService) GetApplicationGroup(ctx context.Context, req *pb.ApplicationGroupIdRequest) (*pb.ApplicationGroup, error) {
	group, err := s.svc.GetApplicationGroup(ctx, req.GroupId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return pbconvert.ToPBApplicationGroup(group), nil
}

func (s *APIService) UpdateApplicationGroup(ctx context.Context, req *pb.UpdateApplicationGroupRequest) (*emptypb.Empty, error) {
	err := s.svc.UpdateApplicationGroup(ctx, req.Id, pbconvert.FromPBUpdateApplicationGroupRequest(req))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *APIService) DeleteApplicationGroup(ctx context.Context, req *pb.ApplicationGroupIdRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteApplicationGroup(ctx, req.GroupId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}
