package repoconvert

import (
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var DeployTypeMapper = mapper.MustNewValueMapper(map[string]domain.DeployType{
	models.ApplicationsDeployTypeRuntime: domain.DeployTypeRuntime,
	models.ApplicationsDeployTypeStatic:  domain.DeployTypeStatic,
})

var ContainerStateMapper = mapper.MustNewValueMapper(map[string]domain.ContainerState{
	models.ApplicationsContainerMissing:  domain.ContainerStateMissing,
	models.ApplicationsContainerStarting: domain.ContainerStateStarting,
	models.ApplicationsContainerRunning:  domain.ContainerStateRunning,
	models.ApplicationsContainerExited:   domain.ContainerStateExited,
	models.ApplicationsContainerErrored:  domain.ContainerStateErrored,
	models.ApplicationsContainerUnknown:  domain.ContainerStateUnknown,
})

func FromDomainApplication(app *domain.Application) *models.Application {
	return &models.Application{
		ID:            app.ID,
		Name:          app.Name,
		RepositoryID:  app.RepositoryID,
		RefName:       app.RefName,
		DeployType:    DeployTypeMapper.FromMust(app.DeployType),
		Running:       app.Running,
		Container:     ContainerStateMapper.FromMust(app.Container),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,
	}
}

func ToDomainApplication(app *models.Application) *domain.Application {
	return &domain.Application{
		ID:            app.ID,
		Name:          app.Name,
		RepositoryID:  app.RepositoryID,
		RefName:       app.RefName,
		DeployType:    DeployTypeMapper.IntoMust(app.DeployType),
		Running:       app.Running,
		Container:     ContainerStateMapper.IntoMust(app.Container),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,

		Config:   ToDomainApplicationConfig(app.R.ApplicationConfig),
		Websites: lo.Map(app.R.Websites, func(website *models.Website, i int) *domain.Website { return ToDomainWebsite(website) }),
		OwnerIDs: lo.Map(app.R.Users, func(user *models.User, i int) string { return user.ID }),
	}
}