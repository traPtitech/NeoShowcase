package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

type BuilderService struct {
	svc usecase.BuilderService

	pb.UnimplementedBuilderServiceServer
}

func NewBuilderServiceServer(svc usecase.BuilderService) *BuilderService {
	return &BuilderService{svc: svc}
}

func (s *BuilderService) GetStatus(ctx context.Context, empty *emptypb.Empty) (*pb.GetStatusResponse, error) {
	return &pb.GetStatusResponse{Status: convertStateToPB(s.svc.GetStatus())}, nil
}

func (s *BuilderService) ConnectEventStream(empty *emptypb.Empty, stream pb.BuilderService_ConnectEventStreamServer) error {
	if err := stream.Send(&pb.Event{Type: pb.Event_CONNECTED, Body: util.ToJSON(map[string]interface{}{})}); err != nil {
		return err
	}

	sub := s.svc.StreamEvents()
	defer sub.Unsubscribe()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case ev := <-sub.Chan():
			evType := ev.Type
			task := ev.Body["task"].(*builder.Task)

			switch evType {
			case event.BuilderBuildStarted:
				if err := stream.Send(&pb.Event{
					Type: pb.Event_BUILD_STARTED,
					Body: util.ToJSON(map[string]interface{}{
						"application_id": task.ApplicationID,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			case event.BuilderBuildFailed:
				if err := stream.Send(&pb.Event{
					Type: pb.Event_BUILD_FAILED,
					Body: util.ToJSON(map[string]interface{}{
						"application_id": task.ApplicationID,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			case event.BuilderBuildSucceeded:
				if err := stream.Send(&pb.Event{
					Type: pb.Event_BUILD_SUCCEEDED,
					Body: util.ToJSON(map[string]interface{}{
						"application_id": task.ApplicationID,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}

			case event.BuilderBuildCanceled:
				if err := stream.Send(&pb.Event{
					Type: pb.Event_BUILD_CANCELED,
					Body: util.ToJSON(map[string]interface{}{
						"application_id": task.ApplicationID,
						"build_id":       task.BuildID,
					}),
				}); err != nil {
					return err
				}
			}
		}
	}
}

func (s *BuilderService) StartBuildImage(ctx context.Context, request *pb.StartBuildImageRequest) (*pb.StartBuildImageResponse, error) {
	task := &builder.Task{
		Static:        false,
		BuildSource:   convertBuildSourceFromPB(request.Source),
		BuildOptions:  convertBuildOptionsFromPB(request.Options),
		ImageName:     request.ImageName,
		ImageTag:      request.ImageTag,
		BuildID:       request.BuildId,
		ApplicationID: request.ApplicationId,
	}

	err := s.svc.StartBuild(ctx, task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.StartBuildImageResponse{}, nil
}

func (s *BuilderService) StartBuildStatic(ctx context.Context, request *pb.StartBuildStaticRequest) (*pb.StartBuildStaticResponse, error) {
	task := &builder.Task{
		Static:        true,
		BuildSource:   convertBuildSourceFromPB(request.Source),
		BuildOptions:  convertBuildOptionsFromPB(request.Options),
		BuildID:       request.BuildId,
		ApplicationID: request.ApplicationId,
	}

	err := s.svc.StartBuild(ctx, task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.StartBuildStaticResponse{}, nil
}

func (s *BuilderService) CancelTask(ctx context.Context, _ *emptypb.Empty) (*pb.CancelTaskResponse, error) {
	id, err := s.svc.CancelBuild(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	if len(id) == 0 {
		return &pb.CancelTaskResponse{Canceled: false}, nil
	}
	return &pb.CancelTaskResponse{Canceled: true, BuildId: id}, nil
}

func convertStateToPB(state builder.State) pb.BuilderStatus {
	switch state {
	case builder.StateUnavailable:
		return pb.BuilderStatus_UNAVAILABLE
	case builder.StateWaiting:
		return pb.BuilderStatus_WAITING
	case builder.StateBuilding:
		return pb.BuilderStatus_BUILDING
	default:
		return pb.BuilderStatus_UNKNOWN
	}
}

func convertBuildSourceFromPB(source *pb.BuildSource) *builder.BuildSource {
	if source == nil {
		return nil
	}
	return &builder.BuildSource{
		RepositoryUrl: source.RepositoryUrl,
		Commit:        source.Commit,
	}
}

func convertBuildOptionsFromPB(options *pb.BuildOptions) *builder.BuildOptions {
	if options == nil {
		return nil
	}
	return &builder.BuildOptions{
		BaseImageName: options.BaseImageName,
		Workdir:       options.Workdir,
		ArtifactPath:  options.ArtifactPath,
		BuildCmd:      options.BuildCmd,
		EntrypointCmd: options.EntrypointCmd,
	}
}
