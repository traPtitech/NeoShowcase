package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
)

func FromDomainEnvironment(env *domain.Environment) *models.Environment {
	return &models.Environment{
		ApplicationID: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}

func ToDomainEnvironment(env *models.Environment) *domain.Environment {
	return &domain.Environment{
		ApplicationID: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}
