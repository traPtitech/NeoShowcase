package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var AuthTypeMapper = mapper.MustNewValueMapper(map[string]domain.AuthenticationType{
	models.WebsitesAuthenticationOff:  domain.AuthenticationTypeOff,
	models.WebsitesAuthenticationSoft: domain.AuthenticationTypeSoft,
	models.WebsitesAuthenticationHard: domain.AuthenticationTypeHard,
})

func FromDomainWebsite(appID string, website *domain.Website) *models.Website {
	return &models.Website{
		ID:             website.ID,
		FQDN:           website.FQDN,
		PathPrefix:     website.PathPrefix,
		StripPrefix:    website.StripPrefix,
		HTTPS:          website.HTTPS,
		H2C:            website.H2C,
		HTTPPort:       website.HTTPPort,
		Authentication: AuthTypeMapper.FromMust(website.Authentication),
		ApplicationID:  appID,
	}
}

func ToDomainWebsite(website *models.Website) *domain.Website {
	return &domain.Website{
		ID:             website.ID,
		FQDN:           website.FQDN,
		PathPrefix:     website.PathPrefix,
		StripPrefix:    website.StripPrefix,
		HTTPS:          website.HTTPS,
		H2C:            website.H2C,
		HTTPPort:       website.HTTPPort,
		Authentication: AuthTypeMapper.IntoMust(website.Authentication),
	}
}
