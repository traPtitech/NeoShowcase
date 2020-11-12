package builder

import (
	"context"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/idgen"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (s *Service) GetStatus(_ context.Context, _ *emptypb.Empty) (*api.GetStatusResponse, error) {
	s.stateLock.RLock()
	defer s.stateLock.RUnlock()
	return &api.GetStatusResponse{Status: s.state}, nil
}

func (s *Service) ConnectEventStream(_ *emptypb.Empty, stream api.BuilderService_ConnectEventStreamServer) error {
	sub := s.bus.Subscribe(10, "builder.*")

	if err := stream.Send(&api.Event{Type: api.Event_CONNECTED}); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			go func() {
				for range sub.Receiver {
				}
			}()
			s.bus.Unsubscribe(sub)
			return nil
		case ev := <-sub.Receiver:
			switch ev.Name {
			case IEventBuildStarted:

			case IEventBuildSucceeded:

			case IEventBuildFailed:

			case IEventBuildCanceled:

			}
		}
	}
}

func (s *Service) StartBuildImageTask(ctx context.Context, req *api.StartBuildImageTaskRequest) (*api.StartBuildImageTaskResponse, error) {
	s.stateLock.Lock()
	if s.state != api.BuilderStatus_WAITING {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Unavailable, "status: %s", s.state.String())
	}

	t := &Task{
		BuildID:       idgen.New(),
		RepositoryURL: req.GetRepositoryUrl(),
		ImageName:     req.GetImageName(),
		BuildLogM: models.BuildLog{
			ApplicationID: req.GetApplicationId(),
			Result:        models.BuildLogsResultBUILDING,
			StartedAt:     time.Now(),
		},
	}
	if err := t.startAsync(ctx, s); err != nil {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	s.state = api.BuilderStatus_BUILDING
	s.stateLock.Unlock()
	return &api.StartBuildImageTaskResponse{BuildId: t.BuildID}, nil
}

func (s *Service) StartBuildStaticTask(ctx context.Context, req *api.StartBuildStaticTaskRequest) (*api.StartBuildStaticTaskResponse, error) {
	s.stateLock.Lock()
	if s.state != api.BuilderStatus_WAITING {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Unavailable, "status: %s", s.state.String())
	}

	t := &Task{
		Static:        true,
		BuildID:       idgen.New(),
		RepositoryURL: req.GetRepositoryUrl(),
		BuildLogM: models.BuildLog{
			ApplicationID: req.GetApplicationId(),
			Result:        models.BuildLogsResultBUILDING,
			StartedAt:     time.Now(),
		},
	}

	if err := t.startAsync(ctx, s); err != nil {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	s.state = api.BuilderStatus_BUILDING
	s.stateLock.Unlock()
	return &api.StartBuildStaticTaskResponse{BuildId: t.BuildID}, nil
}

func (s *Service) CancelTask(_ context.Context, _ *emptypb.Empty) (*api.CancelTaskResponse, error) {
	s.stateLock.RLock()
	state := s.state
	s.stateLock.RUnlock()

	if state != api.BuilderStatus_BUILDING {
		return &api.CancelTaskResponse{Canceled: false}, nil
	}

	s.taskLock.Lock()
	defer s.taskLock.Unlock()
	s.task.cancelFunc()
	return &api.CancelTaskResponse{Canceled: true, BuildId: s.task.BuildID}, nil
}
