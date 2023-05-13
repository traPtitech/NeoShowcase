package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
)

func FromDomainAvailablePort(ap *domain.AvailablePort) *models.AvailablePort {
	return &models.AvailablePort{
		StartPort: ap.StartPort,
		EndPort:   ap.EndPort,
		Protocol:  PortPublicationProtocolMapper.FromMust(ap.Protocol),
	}
}

func ToDomainAvailablePort(ap *models.AvailablePort) *domain.AvailablePort {
	return &domain.AvailablePort{
		StartPort: ap.StartPort,
		EndPort:   ap.EndPort,
		Protocol:  PortPublicationProtocolMapper.IntoMust(ap.Protocol),
	}
}
