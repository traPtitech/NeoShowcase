package usecase

import (
	"context"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIServerService) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	if err := repo.Validate(); err != nil {
		return newError(ErrorTypeBadRequest, "invalid repository", err)
	}

	return s.gitRepo.CreateRepository(ctx, repo)
}

type GetRepoScope int

const (
	GetRepoScopeMine GetRepoScope = iota
	GetRepoScopePublic
	GetRepoScopeAll
)

func (s *APIServerService) GetRepositories(ctx context.Context, scope GetRepoScope) ([]*domain.Repository, error) {
	var cond domain.GetRepositoryCondition
	switch scope {
	case GetRepoScopeMine:
		cond.UserID = optional.From(web.GetUser(ctx).ID)
	case GetRepoScopePublic:
		cond.PublicOrOwnedBy = optional.From(web.GetUser(ctx).ID)
	case GetRepoScopeAll:
		if err := s.isAdmin(ctx); err != nil {
			return nil, err
		}
	}
	return s.gitRepo.GetRepositories(ctx, cond)
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

func (s *APIServerService) RefreshRepository(ctx context.Context, id string) error {
	err := s.controller.FetchRepository(ctx, id)
	if err != nil {
		return errors.Wrap(err, "requesting controller")
	}
	return nil
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
