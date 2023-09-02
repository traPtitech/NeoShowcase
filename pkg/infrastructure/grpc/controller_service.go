package grpc

import (
	"context"

	"connectrpc.com/connect"
	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
	backend    domain.Backend
	fetcher    repofetcher.Service
	cd         cdservice.Service
	builder    domain.ControllerBuilderService
	logStream  *logstream.Service
	pubKey     *ssh.PublicKeys
	sshConf    domain.SSHConfig
	adminerURL string
}

func NewControllerService(
	backend domain.Backend,
	fetcher repofetcher.Service,
	cd cdservice.Service,
	builder domain.ControllerBuilderService,
	logStream *logstream.Service,
	pubKey *ssh.PublicKeys,
	sshConf domain.SSHConfig,
	adminerURL domain.AdminerURL,
) pbconnect.ControllerServiceHandler {
	return &ControllerService{
		backend:    backend,
		fetcher:    fetcher,
		cd:         cd,
		builder:    builder,
		logStream:  logStream,
		pubKey:     pubKey,
		sshConf:    sshConf,
		adminerURL: string(adminerURL),
	}
}

func (s *ControllerService) GetSystemInfo(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.SystemInfo], error) {
	domains := s.backend.AvailableDomains()
	ports := s.backend.AvailablePorts()
	res := connect.NewResponse(&pb.SystemInfo{
		PublicKey: domain.Base64EncodedPublicKey(s.pubKey.Signer.PublicKey()) + " neoshowcase",
		Ssh: &pb.SSHInfo{
			Host: s.sshConf.Host,
			Port: int32(s.sshConf.Port),
		},
		Domains:    ds.Map(domains, pbconvert.ToPBAvailableDomain),
		Ports:      ds.Map(ports, pbconvert.ToPBAvailablePort),
		AdminerUrl: s.adminerURL,
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
