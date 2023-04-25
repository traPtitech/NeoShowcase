package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

func FromDomainUser(user *domain.User) *models.User {
	return &models.User{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}

func ToDomainUser(user *models.User) *domain.User {
	return &domain.User{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}
