package usecase

import (
	"context"
	"fmt"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type CreateRepositoryAuth struct {
	Method   domain.RepositoryAuthMethod
	Username string
	Password string
	KeyID    string
}

func (s *APIServerService) convertRepositoryAuth(a CreateRepositoryAuth) (domain.RepositoryAuth, error) {
	switch a.Method {
	case domain.RepositoryAuthMethodBasic:
		return domain.RepositoryAuth{
			Method:   domain.RepositoryAuthMethodBasic,
			Username: a.Username,
			Password: a.Password,
		}, nil
	case domain.RepositoryAuthMethodSSH:
		key, ok := s.tmpKeys.GetIfExists(a.KeyID)
		if !ok {
			return domain.RepositoryAuth{}, newError(ErrorTypeBadRequest, fmt.Sprintf("key %v does not exist", a.KeyID), nil)
		}
		s.tmpKeys.Forget(a.KeyID)
		pem, err := domain.EncodePrivateKeyPem(key)
		if err != nil {
			return domain.RepositoryAuth{}, err
		}
		return domain.RepositoryAuth{
			Method: domain.RepositoryAuthMethodSSH,
			SSHKey: pem,
		}, nil
	default:
		return domain.RepositoryAuth{}, errors.Errorf("unknown auth method: %v", a.Method)
	}
}

func (s *APIServerService) CreateRepository(ctx context.Context, name, url string, auth optional.Of[CreateRepositoryAuth]) (*domain.Repository, error) {
	dAuth, err := optional.MapErr(auth, s.convertRepositoryAuth)
	if err != nil {
		return nil, err
	}
	user := web.GetUser(ctx)
	repo := domain.NewRepository(name, url, dAuth, []string{user.ID})

	if err := repo.Validate(); err != nil {
		return nil, newError(ErrorTypeBadRequest, "invalid repository", err)
	}

	return repo, s.gitRepo.CreateRepository(ctx, repo)
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

type UpdateRepositoryArgs struct {
	Name     optional.Of[string]
	URL      optional.Of[string]
	Auth     optional.Of[optional.Of[CreateRepositoryAuth]]
	OwnerIDs optional.Of[[]string]
}

func (s *APIServerService) convertUpdateRepositoryArgs(a *UpdateRepositoryArgs) (*domain.UpdateRepositoryArgs, error) {
	dAuth, err := optional.MapErr(a.Auth, func(t optional.Of[CreateRepositoryAuth]) (optional.Of[domain.RepositoryAuth], error) {
		return optional.MapErr(t, s.convertRepositoryAuth)
	})
	if err != nil {
		return nil, err
	}
	return &domain.UpdateRepositoryArgs{
		Name:     a.Name,
		URL:      a.URL,
		Auth:     dAuth,
		OwnerIDs: a.OwnerIDs,
	}, nil
}

func (s *APIServerService) UpdateRepository(ctx context.Context, id string, args *UpdateRepositoryArgs) error {
	err := s.isRepositoryOwner(ctx, id)
	if err != nil {
		return err
	}

	dArgs, err := s.convertUpdateRepositoryArgs(args)
	if err != nil {
		return newError(ErrorTypeBadRequest, "invalid args", err)
	}

	repo, err := s.gitRepo.GetRepository(ctx, id)
	if err != nil {
		return err
	}
	repo.Apply(dArgs)
	if err = repo.Validate(); err != nil {
		return newError(ErrorTypeBadRequest, "invalid repository", err)
	}

	return s.gitRepo.UpdateRepository(ctx, id, dArgs)
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
