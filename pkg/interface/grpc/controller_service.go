package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type ControllerService struct {
	backend   domain.Backend
	fetcher   usecase.RepositoryFetcherService
	cd        usecase.ContinuousDeploymentService
	builder   domain.ControllerBuilderService
	logStream *usecase.LogStreamService
}

func NewControllerService(
	backend domain.Backend,
	fetcher usecase.RepositoryFetcherService,
	cd usecase.ContinuousDeploymentService,
	builder domain.ControllerBuilderService,
	logStream *usecase.LogStreamService,
) pbconnect.ControllerServiceHandler {
	return &ControllerService{
		backend:   backend,
		fetcher:   fetcher,
		cd:        cd,
		builder:   builder,
		logStream: logStream,
	}
}

func (s *ControllerService) AuthAvailable(_ context.Context, c *connect.Request[pb.AuthAvailableRequest]) (*connect.Response[pb.AuthAvailableResponse], error) {
	available := s.backend.AuthAllowed(c.Msg.Fqdn)
	res := connect.NewResponse(&pb.AuthAvailableResponse{Available: available})
	return res, nil
}

func (s *ControllerService) FetchRepository(_ context.Context, c *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[emptypb.Empty], error) {
	s.fetcher.Fetch(c.Msg.RepositoryId)
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) RegisterBuilds(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	s.cd.RegisterBuilds()
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) SyncDeployments(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	s.cd.SyncDeployments()
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) StreamBuildLog(ctx context.Context, c *connect.Request[pb.BuildIdRequest], c2 *connect.ServerStream[pb.BuildLog]) error {
	sub := make(chan []byte, 100)
	ok, unsubscribe := s.logStream.SubscribeBuildLog(c.Msg.BuildId, sub)
	if !ok {
		return errors.New("build log stream unavailable")
	}
	defer unsubscribe()

loop:
	for {
		select {
		case l, ok := <-sub:
			if !ok {
				break loop
			}
			err := c2.Send(&pb.BuildLog{Log: l})
			if err != nil {
				return errors.New("failed to send message")
			}
		case <-ctx.Done():
			break loop
		}
	}
	return nil
}

func (s *ControllerService) CancelBuild(_ context.Context, c *connect.Request[pb.BuildIdRequest]) (*connect.Response[emptypb.Empty], error) {
	buildID := c.Msg.BuildId
	s.builder.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_CANCEL_BUILD,
		Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.BuildIdRequest{BuildId: buildID}},
	})
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
