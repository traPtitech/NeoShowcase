package grpc

import (
	"connectrpc.com/connect"
	"context"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"io"
	"sync"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type builderConnection struct {
	reqSender chan<- *pb.BuilderRequest
	priority  int64
	buildID   string
}

func (c *builderConnection) Send(req *pb.BuilderRequest) {
	select {
	case c.reqSender <- req:
	default:
	}
}

func (c *builderConnection) SetBuildID(id string) {
	c.buildID = id
}

func (c *builderConnection) ClearBuildID() {
	c.buildID = ""
}

func (c *builderConnection) Busy() bool {
	return c.buildID != ""
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
			case pb.BuilderResponse_CONNECTED:
				payload := res.Body.(*pb.BuilderResponse_Connected).Connected
				conn.priority = payload.Priority
			case pb.BuilderResponse_BUILD_STARTED:
				payload := res.Body.(*pb.BuilderResponse_Started).Started
				conn.SetBuildID(payload.BuildId)
				s.logStream.StartBuildLog(payload.BuildId)
			case pb.BuilderResponse_BUILD_SETTLED:
				payload := res.Body.(*pb.BuilderResponse_Settled).Settled
				conn.ClearBuildID()
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

func (s *ControllerBuilderService) StartBuilds(buildIDs []string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Select available builders (and copy the slice)
	conns := lo.Filter(s.builderConnections, func(c *builderConnection, _ int) bool { return !c.Busy() })
	// Select from higher priority builders
	slices.SortFunc(conns, ds.MoreFunc(func(c *builderConnection) int64 { return c.priority }))

	// Send builds to available builders
	for i, conn := range ds.FirstN(conns, len(buildIDs)) {
		buildID := buildIDs[i]
		conn.Send(&pb.BuilderRequest{
			Type: pb.BuilderRequest_START_BUILD,
			Body: &pb.BuilderRequest_StartBuild{StartBuild: &pb.StartBuildRequest{
				BuildId: buildID,
			}},
		})
	}
}

func (s *ControllerBuilderService) CancelBuild(buildID string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	conns := lo.Filter(s.builderConnections, func(c *builderConnection, _ int) bool { return c.buildID == buildID })
	// assert len(conns) <= 1
	for _, conn := range conns {
		conn.Send(&pb.BuilderRequest{
			Type: pb.BuilderRequest_CANCEL_BUILD,
			Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.BuildIdRequest{
				BuildId: buildID,
			}},
		})
	}
}
