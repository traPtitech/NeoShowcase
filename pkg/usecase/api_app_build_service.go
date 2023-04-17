package usecase

import (
	"context"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIServerService) GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(applicationID)})
}

func (s *APIServerService) GetBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	return handleRepoError(build, err)
}

func (s *APIServerService) RetryCommitBuild(ctx context.Context, applicationID string, commit string) error {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return err
	}

	err = s.buildRepo.MarkCommitAsRetriable(ctx, applicationID, commit)
	if err != nil {
		return err
	}
	// NOTE: requires the app to be running for builds to register
	s.bus.Publish(event.CDServiceRegisterBuildRequest, nil)
	return nil
}

func (s *APIServerService) CancelBuild(ctx context.Context, buildID string) error {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return err
	}

	s.component.BroadcastBuilder(&pb.BuilderRequest{
		Type: pb.BuilderRequest_CANCEL_BUILD,
		Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.BuildIdRequest{BuildId: buildID}},
	})
	return nil
}

func (s *APIServerService) GetBuildLog(ctx context.Context, buildID string) ([]byte, error) {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return nil, err
	}

	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		return nil, err
	}
	if !build.Status.IsFinished() {
		return nil, newError(ErrorTypeBadRequest, "build not finished", nil)
	}
	return domain.GetBuildLog(s.storage, buildID)
}

func (s *APIServerService) GetBuildLogStream(ctx context.Context, buildID string, send func(b []byte) error) error {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return err
	}

	sub := make(chan []byte, 100)
	ok, unsubscribe := s.logSvc.SubscribeBuildLog(buildID, sub)
	if !ok {
		return newError(ErrorTypeBadRequest, "build log stream not available", nil)
	}
	defer unsubscribe()

	for b := range sub {
		err = send(b)
		if err != nil {
			return errors.Wrap(err, "failed to send log")
		}
	}
	return nil
}

func (s *APIServerService) GetArtifact(_ context.Context, artifactID string) ([]byte, error) {
	return domain.GetArtifact(s.storage, artifactID)
}
