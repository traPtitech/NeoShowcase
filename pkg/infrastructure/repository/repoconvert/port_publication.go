package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var PortPublicationProtocolMapper = mapper.MustNewValueMapper(map[string]domain.PortPublicationProtocol{
	models.PortPublicationsProtocolTCP: domain.PortPublicationProtocolTCP,
	models.PortPublicationsProtocolUDP: domain.PortPublicationProtocolUDP,
})

func FromDomainPortPublication(appID string, p *domain.PortPublication) *models.PortPublication {
	return &models.PortPublication{
		ApplicationID:   appID,
		InternetPort:    p.InternetPort,
		ApplicationPort: p.ApplicationPort,
		Protocol:        PortPublicationProtocolMapper.FromMust(p.Protocol),
	}
}

func ToDomainPortPublication(p *models.PortPublication) *domain.PortPublication {
	return &domain.PortPublication{
		InternetPort:    p.InternetPort,
		ApplicationPort: p.ApplicationPort,
		Protocol:        PortPublicationProtocolMapper.IntoMust(p.Protocol),
	}
}
