package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIService) CreateRepository(ctx context.Context, req *connect.Request[pb.CreateRepositoryRequest]) (*connect.Response[pb.Repository], error) {
	msg := req.Msg
	user := web.GetUser(ctx)
	repo := &domain.Repository{
		ID:       domain.NewID(),
		Name:     msg.Name,
		URL:      msg.Url,
		Auth:     pbconvert.FromPBRepositoryAuth(msg.Auth),
		OwnerIDs: []string{user.ID},
	}
	err := s.svc.CreateRepository(ctx, repo)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBRepository(repo))
	return res, nil
}

func (s *APIService) GetRepositories(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetRepositoriesResponse], error) {
	repositories, err := s.svc.GetRepositories(ctx)
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

func (s *APIService) UpdateRepository(ctx context.Context, req *connect.Request[pb.UpdateRepositoryRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	args := &domain.UpdateRepositoryArgs{
		Name:     optional.From(msg.Name),
		URL:      optional.From(msg.Url),
		Auth:     optional.From(pbconvert.FromPBRepositoryAuth(msg.Auth)),
		OwnerIDs: optional.From(msg.OwnerIds),
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
