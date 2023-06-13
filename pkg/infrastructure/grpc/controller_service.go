package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type ControllerService struct {
	backend   domain.Backend
	fetcher   repofetcher.Service
	cd        cdservice.Service
	builder   domain.ControllerBuilderService
	logStream *logstream.Service
	sshConf   domain.SSHConfig
}

func NewControllerService(
	backend domain.Backend,
	fetcher repofetcher.Service,
	cd cdservice.Service,
	builder domain.ControllerBuilderService,
	logStream *logstream.Service,
	sshConf domain.SSHConfig,
) pbconnect.ControllerServiceHandler {
	return &ControllerService{
		backend:   backend,
		fetcher:   fetcher,
		cd:        cd,
		builder:   builder,
		logStream: logStream,
		sshConf:   sshConf,
	}
}

func (s *ControllerService) GetSSHInfo(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.SSHInfo], error) {
	res := connect.NewResponse(&pb.SSHInfo{
		Host: s.sshConf.Host,
		Port: int32(s.sshConf.Port),
	})
	return res, nil
}

func (s *ControllerService) GetAvailableDomains(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailableDomains], error) {
	ad := s.backend.AvailableDomains()
	res := connect.NewResponse(&pb.AvailableDomains{
		Domains: ds.Map(ad, pbconvert.ToPBAvailableDomain),
	})
	return res, nil
}

func (s *ControllerService) GetAvailablePorts(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailablePorts], error) {
	ap := s.backend.AvailablePorts()
	res := connect.NewResponse(&pb.AvailablePorts{
		AvailablePorts: ds.Map(ap, pbconvert.ToPBAvailablePort),
	})
	return res, nil
}

func (s *ControllerService) FetchRepository(_ context.Context, c *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[emptypb.Empty], error) {
	s.fetcher.Fetch(c.Msg.RepositoryId)
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) RegisterBuild(_ context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	s.cd.RegisterBuild(req.Msg.Id)
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
