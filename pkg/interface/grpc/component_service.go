package grpc

import (
	"context"
	"io"
	"sync"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

type ComponentServiceGRPCServer struct {
	*grpc.Server
}

func NewComponentServiceGRPCServer() ComponentServiceGRPCServer {
	return ComponentServiceGRPCServer{NewServer()}
}

type builderConnection struct {
	req chan<- *pb.BuilderRequest
}

type ssGenConnection struct {
	req chan<- *pb.SSGenRequest
}

type ComponentService struct {
	bus domain.Bus

	builderConnections []*builderConnection
	ssGenConnections   []*ssGenConnection
	lock               sync.Mutex

	pb.UnimplementedComponentServiceServer
}

func NewComponentServiceServer(bus domain.Bus) domain.ComponentService {
	return &ComponentService{
		bus: bus,
	}
}

func (s *ComponentService) ConnectBuilder(c pb.ComponentService_ConnectBuilderServer) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new builder connection")
	defer log.WithField("id", id).Info("builder connection closed")

	ctx, cancel := context.WithCancel(c.Context())
	req := make(chan *pb.BuilderRequest)
	conn := &builderConnection{req: req}
	s.lock.Lock()
	s.builderConnections = append(s.builderConnections, conn)
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.builderConnections = lo.Without(s.builderConnections, conn)
	}()

	go func() {
		defer cancel()

		for {
			res, err := c.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.WithError(err).Errorf("error receiving builder event: %v", err)
				return
			}

			payload, err := util.IntoField(res.Body)
			if err != nil {
				log.WithError(err).Errorf("failed to decode builder event body: %v", err)
			}

			s.lock.Lock()
			switch res.Type {
			case pb.BuilderResponse_BUILD_STARTED:
				s.bus.Publish(event.BuilderBuildStarted, payload)
			case pb.BuilderResponse_BUILD_SETTLED:
				s.bus.Publish(event.BuilderBuildSettled, payload)
			}
			s.lock.Unlock()
		}
	}()

loop:
	for {
		select {
		case req := <-req:
			err := c.Send(req)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			break loop
		}
	}

	return nil
}

func (s *ComponentService) ConnectSSGen(_ *emptypb.Empty, c pb.ComponentService_ConnectSSGenServer) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new ssgen connection")
	defer log.WithField("id", id).Info("ssgen connection closed")

	req := make(chan *pb.SSGenRequest)
	conn := &ssGenConnection{req: req}
	s.lock.Lock()
	s.ssGenConnections = append(s.ssGenConnections, conn)
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.ssGenConnections = lo.Without(s.ssGenConnections, conn)
	}()

loop:
	for {
		select {
		case req, ok := <-req:
			if !ok {
				break loop
			}
			err := c.Send(req)
			if err != nil {
				return err
			}
		case <-c.Context().Done():
			break loop
		}
	}

	return nil
}

func (s *ComponentService) BroadcastBuilder(req *pb.BuilderRequest) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, builder := range s.builderConnections {
		select {
		case builder.req <- req:
		default:
		}
	}
}

func (s *ComponentService) BroadcastSSGen(req *pb.SSGenRequest) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, ssgen := range s.ssGenConnections {
		select {
		case ssgen.req <- req:
		default:
		}
	}
}
