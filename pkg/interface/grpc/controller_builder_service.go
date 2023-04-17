package grpc

import (
	"context"
	"io"
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type builderConnection struct {
	reqSender chan<- *pb.BuilderRequest
}

type ControllerBuilderService struct {
	logStream *usecase.LogStreamService

	idle    chan struct{}
	settled chan struct{}

	builderConnections []*builderConnection
	lock               sync.Mutex
}

func NewControllerBuilderService(
	logStream *usecase.LogStreamService,
) domain.ControllerBuilderService {
	return &ControllerBuilderService{
		logStream: logStream,
		idle:      make(chan struct{}),
		settled:   make(chan struct{}),
	}
}

func (s *ControllerBuilderService) ConnectBuilder(ctx context.Context, st *connect.BidiStream[pb.BuilderResponse, pb.BuilderRequest]) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new builder connection")
	defer log.WithField("id", id).Info("builder connection closed")

	s.sendIdle()

	ctx, cancel := context.WithCancel(ctx)
	reqSender := make(chan *pb.BuilderRequest)
	conn := &builderConnection{reqSender: reqSender}
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
			res, err := st.Receive()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				log.Errorf("error receiving builder event: %+v", err)
				return
			}

			s.lock.Lock()
			switch res.Type {
			case pb.BuilderResponse_BUILD_STARTED:
				payload := res.Body.(*pb.BuilderResponse_Started).Started
				s.logStream.StartBuildLog(payload.BuildId)
			case pb.BuilderResponse_BUILD_SETTLED:
				payload := res.Body.(*pb.BuilderResponse_Settled).Settled
				s.sendIdle()
				s.sendSettled()
				s.logStream.CloseBuildLog(payload.BuildId)
			case pb.BuilderResponse_BUILD_LOG:
				payload := res.Body.(*pb.BuilderResponse_Log).Log
				s.logStream.AppendBuildLog(payload.BuildId, payload.Log)
			}
			s.lock.Unlock()
		}
	}()

loop:
	for {
		select {
		case req := <-reqSender:
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

func (s *ControllerBuilderService) sendIdle() {
	select {
	case s.idle <- struct{}{}:
	default:
	}
}

func (s *ControllerBuilderService) sendSettled() {
	select {
	case s.settled <- struct{}{}:
	default:
	}
}

func (s *ControllerBuilderService) ListenBuilderIdle() <-chan struct{} {
	return s.idle
}

func (s *ControllerBuilderService) ListenBuildSettled() <-chan struct{} {
	return s.settled
}

func (s *ControllerBuilderService) BroadcastBuilder(req *pb.BuilderRequest) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, builder := range s.builderConnections {
		select {
		case builder.reqSender <- req:
		default:
		}
	}
}
