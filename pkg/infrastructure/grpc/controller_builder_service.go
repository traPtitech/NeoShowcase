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
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type builderConnection struct {
	reqSender chan<- *pb.BuilderRequest
}

type ControllerBuilderService struct {
	logStream *logstream.Service

	idle    domain.PubSub[struct{}]
	settled domain.PubSub[struct{}]

	builderConnections []*builderConnection
	lock               sync.Mutex
}

func NewControllerBuilderService(
	logStream *logstream.Service,
) domain.ControllerBuilderService {
	return &ControllerBuilderService{
		logStream: logStream,
	}
}

func (s *ControllerBuilderService) ConnectBuilder(ctx context.Context, st *connect.BidiStream[pb.BuilderResponse, pb.BuilderRequest]) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new builder connection")
	defer log.WithField("id", id).Info("builder connection closed")

	s.idle.Publish(struct{}{})

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
				s.idle.Publish(struct{}{})
				s.settled.Publish(struct{}{})
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

func (s *ControllerBuilderService) ListenBuilderIdle() (sub <-chan struct{}, unsub func()) {
	return s.idle.Subscribe()
}

func (s *ControllerBuilderService) ListenBuildSettled() (sub <-chan struct{}, unsub func()) {
	return s.settled.Subscribe()
}

func (s *ControllerBuilderService) broadcast(req *pb.BuilderRequest) {
	for _, builder := range s.builderConnections {
		select {
		case builder.reqSender <- req:
		default:
		}
	}
}

func (s *ControllerBuilderService) StartBuilds(buildIDs []string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Send at most n (= number of builders) build requests
	n := len(s.builderConnections)
	for _, buildID := range ds.FirstN(buildIDs, n) {
		s.broadcast(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD,
			Body: &pb.BuilderRequest_StartBuild{StartBuild: &pb.StartBuildRequest{BuildId: buildID}},
		})
	}
}

func (s *ControllerBuilderService) BroadcastBuilder(req *pb.BuilderRequest) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.broadcast(req)
}
