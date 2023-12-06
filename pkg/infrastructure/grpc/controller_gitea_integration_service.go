package grpc

import (
	"context"
	"sync"

	"connectrpc.com/connect"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

type giteaIntegrationConnection struct {
	reqSender chan<- *pb.GiteaIntegrationRequest
}

type ControllerGiteaIntegrationService struct {
	connections []*giteaIntegrationConnection
	lock        sync.Mutex
}

func NewControllerGiteaIntegrationService() domain.ControllerGiteaIntegrationService {
	return &ControllerGiteaIntegrationService{}
}

func (s *ControllerGiteaIntegrationService) Connect(ctx context.Context, _ *connect.Request[emptypb.Empty], st *connect.ServerStream[pb.GiteaIntegrationRequest]) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new gitea integration connection")
	defer log.WithField("id", id).Info("gitea integration connection closed")

	reqSender := make(chan *pb.GiteaIntegrationRequest)
	conn := &giteaIntegrationConnection{reqSender: reqSender}
	s.lock.Lock()
	s.connections = append(s.connections, conn)
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.connections = lo.Without(s.connections, conn)
	}()

loop:
	for {
		select {
		case req, ok := <-reqSender:
			if !ok {
				break loop
			}
			err := st.Send(req)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			break loop
		}
	}

	return nil
}

func (s *ControllerGiteaIntegrationService) Broadcast(req *pb.GiteaIntegrationRequest) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, ssgen := range s.connections {
		select {
		case ssgen.reqSender <- req:
		default:
		}
	}
}
