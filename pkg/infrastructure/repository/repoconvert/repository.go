package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

var RepoAuthMethodMapper = mapper.MustNewValueMapper(map[string]domain.RepositoryAuthMethod{
	models.RepositoryAuthMethodBasic: domain.RepositoryAuthMethodBasic,
	models.RepositoryAuthMethodSSH:   domain.RepositoryAuthMethodSSH,
})

func FromDomainRepositoryAuth(repositoryID string, auth *domain.RepositoryAuth) *models.RepositoryAuth {
	return &models.RepositoryAuth{
		RepositoryID: repositoryID,
		Method:       RepoAuthMethodMapper.FromMust(auth.Method),
		Username:     auth.Username,
		Password:     auth.Password,
		SSHKey:       auth.SSHKey,
	}
}

func ToDomainRepositoryAuth(auth *models.RepositoryAuth) domain.RepositoryAuth {
	return domain.RepositoryAuth{
		Method:   RepoAuthMethodMapper.IntoMust(auth.Method),
		Username: auth.Username,
		Password: auth.Password,
		SSHKey:   auth.SSHKey,
	}
}

func FromDomainRepository(repo *domain.Repository) *models.Repository {
	return &models.Repository{
		ID:   repo.ID,
		Name: repo.Name,
		URL:  repo.URL,
	}
}

func ToDomainRepository(repo *models.Repository) *domain.Repository {
	ret := &domain.Repository{
		ID:       repo.ID,
		Name:     repo.Name,
		URL:      repo.URL,
		OwnerIDs: ds.Map(repo.R.Users, func(user *models.User) string { return user.ID }),
	}
	if repo.R.RepositoryAuth != nil {
		ret.Auth = optional.From(ToDomainRepositoryAuth(repo.R.RepositoryAuth))
	}
	return ret
}
