package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func FromDomainApplicationGroup(group *domain.ApplicationGroup) *models.ApplicationGroup {
	return &models.ApplicationGroup{
		ID:       group.ID,
		Name:     group.Name,
		Implicit: group.Implicit,
	}
}

func ToDomainApplicationGroup(group *models.ApplicationGroup) *domain.ApplicationGroup {
	return &domain.ApplicationGroup{
		ID:       group.ID,
		Name:     group.Name,
		Implicit: group.Implicit,
		OwnerIDs: ds.Map(group.R.Users, func(user *models.User) string { return user.ID }),
	}
}
