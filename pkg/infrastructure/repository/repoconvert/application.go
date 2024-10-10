package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var DeployTypeMapper = mapper.MustNewValueMapper(map[string]domain.DeployType{
	models.ApplicationsDeployTypeRuntime: domain.DeployTypeRuntime,
	models.ApplicationsDeployTypeStatic:  domain.DeployTypeStatic,
})

var ContainerStateMapper = mapper.MustNewValueMapper(map[string]domain.ContainerState{
	models.ApplicationsContainerMissing:    domain.ContainerStateMissing,
	models.ApplicationsContainerStarting:   domain.ContainerStateStarting,
	models.ApplicationsContainerRestarting: domain.ContainerStateRestarting,
	models.ApplicationsContainerRunning:    domain.ContainerStateRunning,
	models.ApplicationsContainerIdle:       domain.ContainerStateIdle,
	models.ApplicationsContainerExited:     domain.ContainerStateExited,
	models.ApplicationsContainerErrored:    domain.ContainerStateErrored,
	models.ApplicationsContainerUnknown:    domain.ContainerStateUnknown,
})

func FromDomainApplication(app *domain.Application) *models.Application {
	return &models.Application{
		ID:               app.ID,
		Name:             app.Name,
		RepositoryID:     app.RepositoryID,
		RefName:          app.RefName,
		Commit:           app.Commit,
		DeployType:       DeployTypeMapper.FromMust(app.DeployType),
		Running:          app.Running,
		Container:        ContainerStateMapper.FromMust(app.Container),
		ContainerMessage: app.ContainerMessage,
		CurrentBuild:     app.CurrentBuild,
		CreatedAt:        app.CreatedAt,
		UpdatedAt:        app.UpdatedAt,
	}
}

func ToDomainApplication(app *models.Application) *domain.Application {
	return &domain.Application{
		ID:               app.ID,
		Name:             app.Name,
		RepositoryID:     app.RepositoryID,
		RefName:          app.RefName,
		Commit:           app.Commit,
		DeployType:       DeployTypeMapper.IntoMust(app.DeployType),
		Running:          app.Running,
		Container:        ContainerStateMapper.IntoMust(app.Container),
		ContainerMessage: app.ContainerMessage,
		CurrentBuild:     app.CurrentBuild,
		CreatedAt:        app.CreatedAt,
		UpdatedAt:        app.UpdatedAt,

		Config:           ToDomainApplicationConfig(app.R.ApplicationConfig),
		Websites:         ds.Map(app.R.Websites, ToDomainWebsite),
		PortPublications: ds.Map(app.R.PortPublications, ToDomainPortPublication),
		OwnerIDs:         ds.Map(app.R.Users, func(user *models.User) string { return user.ID }),
	}
}
