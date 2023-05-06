package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
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

func FromDomainUserKey(key *domain.UserKey) *models.UserKey {
	return &models.UserKey{
		ID:        key.ID,
		UserID:    key.UserID,
		PublicKey: key.PublicKey,
	}
}

func ToDomainUserKey(key *models.UserKey) *domain.UserKey {
	return &domain.UserKey{
		ID:        key.ID,
		UserID:    key.UserID,
		PublicKey: key.PublicKey,
	}
}
