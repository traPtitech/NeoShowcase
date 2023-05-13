package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetBuilds(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.GetBuildsResponse], error) {
	builds, err := s.svc.GetBuilds(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetBuildsResponse{
		Builds: ds.Map(builds, pbconvert.ToPBBuild),
	})
	return res, nil
}

func (s *APIService) GetBuild(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.Build], error) {
	build, err := s.svc.GetBuild(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBBuild(build))
	return res, nil
}

func (s *APIService) RetryCommitBuild(ctx context.Context, req *connect.Request[pb.RetryCommitBuildRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	err := s.svc.RetryCommitBuild(ctx, msg.ApplicationId, msg.Commit)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) CancelBuild(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.CancelBuild(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) GetBuildLog(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.BuildLog], error) {
	log, err := s.svc.GetBuildLog(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.BuildLog{Log: log})
	return res, nil
}

func (s *APIService) GetBuildLogStream(ctx context.Context, req *connect.Request[pb.BuildIdRequest], st *connect.ServerStream[pb.BuildLog]) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := s.svc.GetBuildLogStream(ctx, req.Msg.BuildId)
	if err != nil {
		return handleUseCaseError(err)
	}
	for l := range ch {
		err = st.Send(l)
		if err != nil {
			return errors.New("failed to send event")
		}
	}
	return nil
}

func (s *APIService) GetBuildArtifact(ctx context.Context, req *connect.Request[pb.ArtifactIdRequest]) (*connect.Response[pb.ArtifactContent], error) {
	content, err := s.svc.GetArtifact(ctx, req.Msg.ArtifactId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ArtifactContent{
		Filename: req.Msg.ArtifactId + ".tar",
		Content:  content,
	})
	return res, nil
}
