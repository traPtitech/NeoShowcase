package grpc

import (
	"context"
	"io"

	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetAllBuilds(ctx context.Context, c *pb.GetAllBuildsRequest) (*pb.GetBuildsResponse, error) {
	builds, err := s.svc.GetAllBuilds(ctx, int(c.Page), int(c.Limit))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.GetBuildsResponse{
		Builds: ds.Map(builds, pbconvert.ToPBBuild),
	}
	return res, nil
}

func (s *APIService) GetBuilds(ctx context.Context, req *pb.ApplicationIdRequest) (*pb.GetBuildsResponse, error) {
	builds, err := s.svc.GetBuilds(ctx, req.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.GetBuildsResponse{
		Builds: ds.Map(builds, pbconvert.ToPBBuild),
	}
	return res, nil
}

func (s *APIService) GetBuild(ctx context.Context, req *pb.BuildIdRequest) (*pb.Build, error) {
	build, err := s.svc.GetBuild(ctx, req.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := pbconvert.ToPBBuild(build)
	return res, nil
}

func (s *APIService) RetryCommitBuild(ctx context.Context, req *pb.RetryCommitBuildRequest) (*emptypb.Empty, error) {
	msg := req
	err := s.svc.RetryCommitBuild(ctx, msg.ApplicationId, msg.Commit)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) CancelBuild(ctx context.Context, req *pb.BuildIdRequest) (*emptypb.Empty, error) {
	err := s.svc.CancelBuild(ctx, req.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}

func (s *APIService) GetBuildLog(ctx context.Context, req *pb.BuildIdRequest) (*pb.BuildLog, error) {
	log, err := s.svc.GetBuildLog(ctx, req.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.BuildLog{Log: log}
	return res, nil
}

func (s *APIService) GetBuildLogStream(ctx context.Context, req *pb.BuildIdRequest, st pbconnect.APIServiceGetBuildLogStreamServerStream) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := s.svc.GetBuildLogStream(ctx, req.BuildId)
	if err != nil {
		return handleUseCaseError(err)
	}
	for l := range ch {
		err = st.Send(l)
		if err != nil {
			return oops.Wrapf(err, "sending event")
		}
	}
	return nil
}

func (s *APIService) GetBuildArtifact(ctx context.Context, req *pb.ArtifactIdRequest) (*pb.ArtifactContent, error) {
	filename, r, err := s.svc.GetArtifact(ctx, req.ArtifactId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	defer r.Close()
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	res := &pb.ArtifactContent{
		Filename: filename,
		Content:  content,
	}
	return res, nil
}
