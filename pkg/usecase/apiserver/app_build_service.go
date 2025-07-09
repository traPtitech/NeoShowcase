package apiserver

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *Service) GetAllBuilds(ctx context.Context, page, limit int) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{
		Offset:  optional.From(page * limit),
		Limit:   optional.From(limit),
		SortAsc: optional.From(false), // desc
	})
}

func (s *Service) GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(applicationID)})
}

func (s *Service) GetBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	return handleRepoError(build, err)
}

func (s *Service) RetryCommitBuild(ctx context.Context, applicationID string, commit string) error {
	err := s.isApplicationOwner(ctx, applicationID)
	if err != nil {
		return err
	}

	err = s.buildRepo.MarkCommitAsRetriable(ctx, applicationID, commit)
	if err != nil {
		return err
	}
	// NOTE: requires the app to be running for builds to register
	err = s.controller.RegisterBuild(ctx, applicationID)
	if err != nil {
		return errors.Wrap(err, "failed to request new build registration")
	}
	return nil
}

func (s *Service) CancelBuild(ctx context.Context, buildID string) error {
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

func (s *Service) GetBuildLog(ctx context.Context, buildID string) ([]byte, error) {
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

func (s *Service) GetBuildLogStream(ctx context.Context, buildID string) (<-chan *pb.BuildLog, error) {
	err := s.isBuildOwner(ctx, buildID)
	if err != nil {
		return nil, err
	}

	addr, err := s.controller.DiscoverBuildLogInstance(ctx, buildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to discover build log instance")
	}
	if addr.Address == nil {
		return nil, newError(ErrorTypeBadRequest, "build log instance not found", nil)
	}
	ch, err := s.controller.StreamBuildLog(ctx, *addr.Address, buildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to build log stream")
	}
	return ch, nil
}

func (s *Service) GetArtifact(ctx context.Context, artifactID string) (filename string, r io.ReadCloser, err error) {
	artifact, err := s.artifactRepo.GetArtifact(ctx, artifactID)
	if err != nil {
		return "", nil, err
	}
	r, err = domain.GetArtifact(s.storage, artifactID)
	return artifact.Name, r, err
}
