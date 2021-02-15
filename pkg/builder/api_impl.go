package builder

import (
	"context"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/event"
	"github.com/traPtitech/neoshowcase/pkg/idgen"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"github.com/volatiletech/null/v8"
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

	if err := stream.Send(&api.Event{Type: api.Event_CONNECTED, Body: util.ToJSON(map[string]interface{}{})}); err != nil {
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
			case event.BuilderBuildStarted:
				task := ev.Fields["task"].(*Task)
				if err := stream.Send(&api.Event{
					Type: api.Event_BUILD_STARTED,
					Body: util.ToJSON(map[string]interface{}{
						"environment_id": task.BuildLogM.EnvironmentID.String,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			case event.BuilderBuildFailed:
				task := ev.Fields["task"].(*Task)
				if err := stream.Send(&api.Event{
					Type: api.Event_BUILD_FAILED,
					Body: util.ToJSON(map[string]interface{}{
						"environment_id": task.BuildLogM.EnvironmentID.String,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			case event.BuilderBuildSucceeded:
				task := ev.Fields["task"].(*Task)
				if err := stream.Send(&api.Event{
					Type: api.Event_BUILD_SUCCEEDED,
					Body: util.ToJSON(map[string]interface{}{
						"environment_id": task.BuildLogM.EnvironmentID.String,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			case event.BuilderBuildCanceled:
				task := ev.Fields["task"].(*Task)
				if err := stream.Send(&api.Event{
					Type: api.Event_BUILD_CANCELED,
					Body: util.ToJSON(map[string]interface{}{
						"environment_id": task.BuildLogM.EnvironmentID.String,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			}
		}
	}
}

func (s *Service) StartBuildImage(ctx context.Context, req *api.StartBuildImageRequest) (*api.StartBuildImageResponse, error) {
	s.stateLock.Lock()
	if s.state != api.BuilderStatus_WAITING {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Unavailable, "status: %s", s.state)
	}

	t := &Task{
		BuildID:      idgen.New(),
		BuildSource:  req.GetSource(),
		BuildOptions: req.GetOptions(),
		ImageName:    req.GetImageName(),
		BuildLogM: models.BuildLog{
			Result:    models.BuildLogsResultBUILDING,
			StartedAt: time.Now(),
		},
	}
	// 環境IDが指定されていない場合はデバッグビルド
	if len(req.EnvironmentId) > 0 {
		t.BuildLogM.EnvironmentID = null.StringFrom(req.GetEnvironmentId())
	}

	if err := t.startAsync(ctx, s); err != nil {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	s.state = api.BuilderStatus_BUILDING
	s.stateLock.Unlock()
	return &api.StartBuildImageResponse{BuildId: t.BuildID}, nil
}

func (s *Service) StartBuildStatic(ctx context.Context, req *api.StartBuildStaticRequest) (*api.StartBuildStaticResponse, error) {
	s.stateLock.Lock()
	if s.state != api.BuilderStatus_WAITING {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Unavailable, "status: %s", s.state)
	}

	t := &Task{
		Static:       true,
		BuildID:      idgen.New(),
		BuildSource:  req.GetSource(),
		BuildOptions: req.GetOptions(),
		BuildLogM: models.BuildLog{
			Result:    models.BuildLogsResultBUILDING,
			StartedAt: time.Now(),
		},
	}
	// アプリケーションIDが指定されていない場合はデバッグビルド
	if len(req.EnvironmentId) > 0 {
		t.BuildLogM.EnvironmentID = null.StringFrom(req.GetEnvironmentId())
	}

	if err := t.startAsync(ctx, s); err != nil {
		s.stateLock.Unlock()
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	s.state = api.BuilderStatus_BUILDING
	s.stateLock.Unlock()
	return &api.StartBuildStaticResponse{BuildId: t.BuildID}, nil
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
