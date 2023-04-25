package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var AuthTypeMapper = mapper.MustNewValueMapper(map[domain.AuthenticationType]pb.AuthenticationType{
	domain.AuthenticationTypeOff:  pb.AuthenticationType_OFF,
	domain.AuthenticationTypeSoft: pb.AuthenticationType_SOFT,
	domain.AuthenticationTypeHard: pb.AuthenticationType_HARD,
})

func FromPBCreateWebsiteRequest(req *pb.CreateWebsiteRequest) *domain.Website {
	return &domain.Website{
		ID:             domain.NewID(),
		FQDN:           req.Fqdn,
		PathPrefix:     req.PathPrefix,
		StripPrefix:    req.StripPrefix,
		HTTPS:          req.Https,
		HTTPPort:       int(req.HttpPort),
		Authentication: AuthTypeMapper.FromMust(req.Authentication),
	}
}

func ToPBWebsite(website *domain.Website) *pb.Website {
	return &pb.Website{
		Id:             website.ID,
		Fqdn:           website.FQDN,
		PathPrefix:     website.PathPrefix,
		Https:          website.HTTPS,
		HttpPort:       int32(website.HTTPPort),
		Authentication: AuthTypeMapper.IntoMust(website.Authentication),
	}
}
