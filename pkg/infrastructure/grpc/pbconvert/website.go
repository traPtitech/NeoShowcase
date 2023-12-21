package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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
		H2C:            req.H2C,
		HTTPPort:       int(req.HttpPort),
		Authentication: AuthTypeMapper.FromMust(req.Authentication),
	}
}

func FromPBUpdateWebsites(req *pb.UpdateApplicationRequest_UpdateWebsites) []*domain.Website {
	return ds.Map(req.Websites, FromPBCreateWebsiteRequest)
}

func ToPBWebsite(website *domain.Website) *pb.Website {
	return &pb.Website{
		Id:             website.ID,
		Fqdn:           website.FQDN,
		PathPrefix:     website.PathPrefix,
		StripPrefix:    website.StripPrefix,
		Https:          website.HTTPS,
		H2C:            website.H2C,
		HttpPort:       int32(website.HTTPPort),
		Authentication: AuthTypeMapper.IntoMust(website.Authentication),
	}
}

func FromPBWebsite(website *pb.Website) *domain.Website {
	return &domain.Website{
		ID:             website.Id,
		FQDN:           website.Fqdn,
		PathPrefix:     website.PathPrefix,
		StripPrefix:    website.StripPrefix,
		HTTPS:          website.Https,
		H2C:            website.H2C,
		HTTPPort:       int(website.HttpPort),
		Authentication: AuthTypeMapper.FromMust(website.Authentication),
	}
}
