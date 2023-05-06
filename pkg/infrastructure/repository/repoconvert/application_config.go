package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
)

func FromDomainApplicationConfig(appID string, c *domain.ApplicationConfig) *models.ApplicationConfig {
	mc := &models.ApplicationConfig{
		ApplicationID: appID,
		UseMariadb:    c.UseMariaDB,
		UseMongodb:    c.UseMongoDB,
		BuildType:     BuildTypeMapper.FromMust(c.BuildConfig.BuildType()),
		Entrypoint:    c.Entrypoint,
		Command:       c.Command,
	}
	FromDomainBuildConfig(c.BuildConfig, mc)
	return mc
}

func ToDomainApplicationConfig(c *models.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		UseMariaDB:  c.UseMariadb,
		UseMongoDB:  c.UseMongodb,
		BuildConfig: ToDomainBuildConfig(c),
		Entrypoint:  c.Entrypoint,
		Command:     c.Command,
	}
}
