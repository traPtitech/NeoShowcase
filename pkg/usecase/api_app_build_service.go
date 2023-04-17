package usecase

import (
	"context"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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
	err = s.controller.RegisterBuilds(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to request new build registration")
	}
	return nil
}

func (s *APIServerService) CancelBuild(ctx context.Context, buildID string) error {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return err
	}

	err = s.controller.CancelBuild(ctx, buildID)
	if err != nil {
		return errors.Wrap(err, "failed to request cancel build")
	}
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

func (s *APIServerService) GetBuildLogStream(ctx context.Context, buildID string) (<-chan *pb.BuildLog, error) {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return nil, err
	}

	ch, err := s.controller.StreamBuildLog(ctx, buildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to build log stream")
	}
	return ch, nil
}

func (s *APIServerService) GetArtifact(_ context.Context, artifactID string) ([]byte, error) {
	return domain.GetArtifact(s.storage, artifactID)
}
