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

type ssGenConnection struct {
	reqSender chan<- *pb.SSGenRequest
}

type ControllerSSGenService struct {
	ssGenConnections []*ssGenConnection
	lock             sync.Mutex
}

func NewControllerSSGenService() domain.ControllerSSGenService {
	return &ControllerSSGenService{}
}

func (s *ControllerSSGenService) ConnectSSGen(ctx context.Context, _ *connect.Request[emptypb.Empty], st *connect.ServerStream[pb.SSGenRequest]) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new ssgen connection")
	defer log.WithField("id", id).Info("ssgen connection closed")

	reqSender := make(chan *pb.SSGenRequest)
	conn := &ssGenConnection{reqSender: reqSender}
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

func (s *ControllerSSGenService) BroadcastSSGen(req *pb.SSGenRequest) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, ssgen := range s.ssGenConnections {
		select {
		case ssgen.reqSender <- req:
		default:
		}
	}
}
