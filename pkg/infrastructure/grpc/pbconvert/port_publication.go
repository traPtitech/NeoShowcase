package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var PortPublicationProtocolMapper = mapper.MustNewValueMapper(map[domain.PortPublicationProtocol]pb.PortPublicationProtocol{
	domain.PortPublicationProtocolTCP: pb.PortPublicationProtocol_TCP,
	domain.PortPublicationProtocolUDP: pb.PortPublicationProtocol_UDP,
})

func FromPBPortPublication(p *pb.PortPublication) *domain.PortPublication {
	return &domain.PortPublication{
		InternetPort:    int(p.InternetPort),
		ApplicationPort: int(p.ApplicationPort),
		Protocol:        PortPublicationProtocolMapper.FromMust(p.Protocol),
	}
}
