package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func FromPBAvailableDomain(ad *pb.AvailableDomain) *domain.AvailableDomain {
	return &domain.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}

func ToPBAvailableDomain(ad *domain.AvailableDomain) *pb.AvailableDomain {
	return &pb.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}
