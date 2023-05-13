package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func FromPBAvailablePort(ap *pb.AvailablePort) *domain.AvailablePort {
	return &domain.AvailablePort{
		StartPort: int(ap.StartPort),
		EndPort:   int(ap.EndPort),
		Protocol:  PortPublicationProtocolMapper.FromMust(ap.Protocol),
	}
}

func ToPBAvailablePort(ap *domain.AvailablePort) *pb.AvailablePort {
	return &pb.AvailablePort{
		StartPort: int32(ap.StartPort),
		EndPort:   int32(ap.EndPort),
		Protocol:  PortPublicationProtocolMapper.IntoMust(ap.Protocol),
	}
}

func ToPBUnavailablePort(up *domain.UnavailablePort) *pb.UnavailablePort {
	return &pb.UnavailablePort{
		Port:     int32(up.Port),
		Protocol: PortPublicationProtocolMapper.IntoMust(up.Protocol),
	}
}
