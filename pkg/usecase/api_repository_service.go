package usecase

import (
	"context"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIServerService) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	if err := repo.Validate(); err != nil {
		return newError(ErrorTypeBadRequest, "invalid repository", err)
	}

	return s.gitRepo.CreateRepository(ctx, repo)
}

func (s *APIServerService) GetRepositories(ctx context.Context) ([]*domain.Repository, error) {
	return s.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{})
}

func (s *APIServerService) GetRepository(ctx context.Context, id string) (*domain.Repository, error) {
	return handleRepoError(s.gitRepo.GetRepository(ctx, id))
}

func (s *APIServerService) UpdateRepository(ctx context.Context, id string, args *domain.UpdateRepositoryArgs) error {
	err := s.isRepositoryOwner(ctx, id)
	if err != nil {
		return err
	}

	repo, err := s.gitRepo.GetRepository(ctx, id)
	if err != nil {
		return err
	}
	repo.Apply(args)
	if err = repo.Validate(); err != nil {
		return newError(ErrorTypeBadRequest, "invalid repository", err)
	}

	return s.gitRepo.UpdateRepository(ctx, id, args)
}

func (s *APIServerService) DeleteRepository(ctx context.Context, id string) error {
	err := s.isRepositoryOwner(ctx, id)
	if err != nil {
		return err
	}

	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{RepositoryID: optional.From(id)})
	if err != nil {
		return errors.Wrap(err, "failed to get related applications")
	}
	if len(apps) > 0 {
		return newError(ErrorTypeBadRequest, "all related applications must be deleted first", nil)
	}

	return s.gitRepo.DeleteRepository(ctx, id)
}
