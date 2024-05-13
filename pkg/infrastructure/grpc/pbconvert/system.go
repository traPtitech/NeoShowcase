package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func FromPBAdditionalLink(link *pb.AdditionalLink) *domain.AdditionalLink {
	return &domain.AdditionalLink{
		Name: link.Name,
		URL:  link.Url,
	}
}

func ToPBAdditionalLink(link *domain.AdditionalLink) *pb.AdditionalLink {
	return &pb.AdditionalLink{
		Name: link.Name,
		Url:  link.URL,
	}
}

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
		AdditionalLinks:  ds.Map(i.AdditionalLinks, FromPBAdditionalLink),
		Version:          i.Version,
		Revision:         i.Revision,
	}
}

func ToPBSystemInfo(i *domain.SystemInfo) *pb.SystemInfo {
	return &pb.SystemInfo{
		PublicKey: i.PublicKey,
		Ssh: &pb.SSHInfo{
			Host: i.SSHInfo.Host,
			Port: int32(i.SSHInfo.Port),
		},
		Domains:         ds.Map(i.AvailableDomains, ToPBAvailableDomain),
		Ports:           ds.Map(i.AvailablePorts, ToPBAvailablePort),
		AdditionalLinks: ds.Map(i.AdditionalLinks, ToPBAdditionalLink),
		Version:         i.Version,
		Revision:        i.Revision,
	}
}
