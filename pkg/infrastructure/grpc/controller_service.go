package grpc

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/friendsofgo/errors"
	"github.com/motoki317/sc"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/systeminfo"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
)

type ControllerPort int

type ControllerService struct {
	Port      int
	cluster   *discovery.Cluster
	infoSvc   systeminfo.Service
	fetcher   repofetcher.Service
	cd        domain.CDService
	builder   domain.ControllerBuilderService
	logStream *logstream.Service

	clientCache *sc.Cache[string, domain.ControllerServiceClient]
}

func NewControllerService(
	port ControllerPort,
	cluster *discovery.Cluster,
	infoSvc systeminfo.Service,
	fetcher repofetcher.Service,
	cd domain.CDService,
	builder domain.ControllerBuilderService,
	logStream *logstream.Service,
) pbconnect.ControllerServiceHandler {
	return &ControllerService{
		Port:      int(port),
		cluster:   cluster,
		infoSvc:   infoSvc,
		fetcher:   fetcher,
		cd:        cd,
		builder:   builder,
		logStream: logStream,

		clientCache: sc.NewMust(func(ctx context.Context, address string) (domain.ControllerServiceClient, error) {
			return NewControllerServiceClient(ControllerServiceClientConfig{URL: address}), nil
		}, time.Hour, time.Hour),
	}
}

func (s *ControllerService) GetSystemInfo(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.SystemInfo], error) {
	info, err := s.infoSvc.GetSystemInfo()
	if err != nil {
		return nil, err
	}
	pbInfo := pbconvert.ToPBSystemInfo(info)
	res := connect.NewResponse(pbInfo)
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

func (s *ControllerService) SyncDeployments(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	s.cd.SyncDeploymentsLocal()

	// Broadcast to cluster
	p := pool.New().WithErrors()
	for _, address := range s.cluster.AllNeighborAddresses(s.Port) {
		p.Go(func() error {
			client, _ := s.clientCache.Get(ctx, address)
			return client.SyncDeploymentsLocal(ctx)
		})
	}
	err := p.Wait()
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) DiscoverBuildLogInstance(ctx context.Context, c *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.AddressInfo], error) {
	res, err := s.DiscoverBuildLogLocal(ctx, c)
	if err != nil {
		return nil, err
	}
	if res.Msg.Address != nil {
		return res, nil
	}

	// Search for neighbors
	p := pool.NewWithResults[*pb.AddressInfo]().WithErrors()
	for _, address := range s.cluster.AllNeighborAddresses(s.Port) {
		p.Go(func() (*pb.AddressInfo, error) {
			client, _ := s.clientCache.Get(ctx, address)
			return client.DiscoverBuildLogLocal(ctx, c.Msg.BuildId)
		})
	}
	neighborResults, err := p.Wait()
	if err != nil {
		return nil, err
	}
	neighborResult, ok := lo.Find(neighborResults, func(r *pb.AddressInfo) bool {
		return r.Address != nil
	})
	if !ok {
		return nil, errors.New("build log not available")
	}
	return connect.NewResponse(neighborResult), nil
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

func (s *ControllerService) StartBuild(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	s.cd.StartBuildLocal()

	// Broadcast to cluster
	p := pool.New().WithErrors()
	for _, address := range s.cluster.AllNeighborAddresses(s.Port) {
		p.Go(func() error {
			client, _ := s.clientCache.Get(ctx, address)
			return client.StartBuildLocal(ctx)
		})
	}
	err := p.Wait()
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) CancelBuild(ctx context.Context, c *connect.Request[pb.BuildIdRequest]) (*connect.Response[emptypb.Empty], error) {
	buildID := c.Msg.BuildId
	s.builder.CancelBuild(buildID)

	// Broadcast to cluster
	p := pool.New().WithErrors()
	for _, address := range s.cluster.AllNeighborAddresses(s.Port) {
		p.Go(func() error {
			client, _ := s.clientCache.Get(ctx, address)
			return client.CancelBuildLocal(ctx, buildID)
		})
	}
	err := p.Wait()
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) DiscoverBuildLogLocal(_ context.Context, c *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.AddressInfo], error) {
	ok := s.logStream.HasBuildLog(c.Msg.BuildId)
	if !ok {
		return connect.NewResponse(&pb.AddressInfo{}), nil
	}
	addr, ok := s.cluster.MyAddress(s.Port)
	if !ok {
		return nil, errors.New("self address not available")
	}
	return connect.NewResponse(&pb.AddressInfo{Address: &addr}), nil
}

func (s *ControllerService) StartBuildLocal(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	s.cd.StartBuildLocal()
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) SyncDeploymentsLocal(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	s.cd.SyncDeploymentsLocal()
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerService) CancelBuildLocal(_ context.Context, c *connect.Request[pb.BuildIdRequest]) (*connect.Response[emptypb.Empty], error) {
	buildID := c.Msg.BuildId
	s.builder.CancelBuild(buildID)
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
