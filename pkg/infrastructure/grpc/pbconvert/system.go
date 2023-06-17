package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func FromPBSystemInfo(i *pb.SystemInfo) *domain.SystemInfo {
	return &domain.SystemInfo{
		PublicKey: i.PublicKey,
		SSHInfo: struct {
			Host string
			Port int
		}{
			Host: i.Ssh.Host,
			Port: int(i.Ssh.Port),
		},
		AvailableDomains: ds.Map(i.Domains, FromPBAvailableDomain),
		AvailablePorts:   ds.Map(i.Ports, FromPBAvailablePort),
	}
}

func ToPBSystemInfo(i *domain.SystemInfo) *pb.SystemInfo {
	return &pb.SystemInfo{
		PublicKey: i.PublicKey,
		Ssh: &pb.SSHInfo{
			Host: i.SSHInfo.Host,
			Port: int32(i.SSHInfo.Port),
		},
		Domains: ds.Map(i.AvailableDomains, ToPBAvailableDomain),
		Ports:   ds.Map(i.AvailablePorts, ToPBAvailablePort),
	}
}
