package grpc

import (
	"context"

	"connectrpc.com/connect"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIService) CreateRepository(ctx context.Context, req *connect.Request[pb.CreateRepositoryRequest]) (*connect.Response[pb.Repository], error) {
	msg := req.Msg
	repo, err := s.svc.CreateRepository(ctx,
		msg.Name,
		msg.Url,
		pbconvert.FromPBRepositoryAuth(msg.Auth),
	)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBRepository(repo))
	return res, nil
}

func (s *APIService) GetRepositories(ctx context.Context, req *connect.Request[pb.GetRepositoriesRequest]) (*connect.Response[pb.GetRepositoriesResponse], error) {
	repositories, err := s.svc.GetRepositories(ctx, pbconvert.RepoScopeMapper.FromMust(req.Msg.Scope))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetRepositoriesResponse{
		Repositories: ds.Map(repositories, pbconvert.ToPBRepository),
	})
	return res, nil
}

func (s *APIService) GetRepository(ctx context.Context, req *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[pb.Repository], error) {
	repository, err := s.svc.GetRepository(ctx, req.Msg.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBRepository(repository))
	return res, nil
}

func (s *APIService) GetRepositoryRefs(ctx context.Context, req *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[pb.GetRepositoryRefsResponse], error) {
	refs, err := s.svc.GetRepositoryRefs(ctx, req.Msg.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	pbRefs := lo.MapToSlice(refs, func(ref, commit string) *pb.GitRef {
		return &pb.GitRef{
			RefName: ref,
			Commit:  commit,
		}
	})
	slices.SortFunc(pbRefs, ds.LessFunc(func(r *pb.GitRef) string { return r.RefName }))
	res := connect.NewResponse(&pb.GetRepositoryRefsResponse{Refs: pbRefs})
	return res, nil
}

func (s *APIService) UpdateRepository(ctx context.Context, req *connect.Request[pb.UpdateRepositoryRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	args := &apiserver.UpdateRepositoryArgs{
		Name:     optional.FromPtr(msg.Name),
		URL:      optional.FromPtr(msg.Url),
		Auth:     optional.Map(optional.FromNonZero(msg.Auth), pbconvert.FromPBRepositoryAuth),
		OwnerIDs: optional.Map(optional.FromNonZero(msg.OwnerIds), pbconvert.FromPBUpdateRepositoryOwners),
	}
	err := s.svc.UpdateRepository(ctx, msg.Id, args)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) RefreshRepository(ctx context.Context, req *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.RefreshRepository(ctx, req.Msg.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) DeleteRepository(ctx context.Context, req *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.DeleteRepository(ctx, req.Msg.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
