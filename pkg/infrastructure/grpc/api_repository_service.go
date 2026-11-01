package grpc

import (
	"context"
	"slices"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIService) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.Repository, error) {
	repo, err := s.svc.CreateRepository(ctx,
		req.Name,
		req.Url,
		pbconvert.FromPBRepositoryAuth(req.Auth),
	)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := pbconvert.ToPBRepository(repo)
	return res, nil
}

func (s *APIService) GetRepositories(ctx context.Context, req *pb.GetRepositoriesRequest) (*pb.GetRepositoriesResponse, error) {
	repositories, err := s.svc.GetRepositories(ctx, pbconvert.RepoScopeMapper.FromMust(req.Scope))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.GetRepositoriesResponse{
		Repositories: ds.Map(repositories, pbconvert.ToPBRepository),
	}
	return res, nil
}

func (s *APIService) GetRepositoryCommits(ctx context.Context, req *pb.GetRepositoryCommitsRequest) (*pb.GetRepositoryCommitsResponse, error) {
	commits, err := s.svc.GetRepositoryCommits(ctx, req.Hashes)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.GetRepositoryCommitsResponse{
		Commits: ds.Map(commits, pbconvert.ToPBSimpleCommit),
	}
	return res, nil
}

func (s *APIService) GetRepository(ctx context.Context, req *pb.RepositoryIdRequest) (*pb.Repository, error) {
	repository, err := s.svc.GetRepository(ctx, req.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := pbconvert.ToPBRepository(repository)
	return res, nil
}

func (s *APIService) GetRepositoryRefs(ctx context.Context, req *pb.RepositoryIdRequest) (*pb.GetRepositoryRefsResponse, error) {
	refs, err := s.svc.GetRepositoryRefs(ctx, req.RepositoryId)
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
	res := &pb.GetRepositoryRefsResponse{Refs: pbRefs}
	return res, nil
}

func (s *APIService) UpdateRepository(ctx context.Context, req *pb.UpdateRepositoryRequest) (*emptypb.Empty, error) {
	args := &apiserver.UpdateRepositoryArgs{
		Name:     optional.FromPtr(req.Name),
		URL:      optional.FromPtr(req.Url),
		Auth:     optional.FromNonZero(req.Auth).Map(pbconvert.FromPBRepositoryAuth),
		OwnerIDs: optional.FromNonZero(req.OwnerIds).Map(pbconvert.FromPBUpdateRepositoryOwners),
	}
	err := s.svc.UpdateRepository(ctx, req.Id, args)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) RefreshRepository(ctx context.Context, req *pb.RepositoryIdRequest) (*emptypb.Empty, error) {
	err := s.svc.RefreshRepository(ctx, req.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) DeleteRepository(ctx context.Context, req *pb.RepositoryIdRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteRepository(ctx, req.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}
