package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var BuildTypeMapper = mapper.NewValueMapper(map[domain.BuildType]pb.BuildType{
	domain.BuildTypeRuntimeCmd:        pb.BuildType_RUNTIME_CMD,
	domain.BuildTypeRuntimeDockerfile: pb.BuildType_RUNTIME_DOCKERFILE,
	domain.BuildTypeStaticCmd:         pb.BuildType_STATIC_CMD,
	domain.BuildTypeStaticDockerfile:  pb.BuildType_STATIC_DOCKERFILE,
})

func FromPBApplicationConfig(c *pb.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		UseMariaDB:  c.UseMariadb,
		UseMongoDB:  c.UseMongodb,
		BuildType:   BuildTypeMapper.FromMust(c.BuildType),
		BuildConfig: FromPBBuildConfig(c.BuildConfig),
	}
}

func ToPBApplicationConfig(c domain.ApplicationConfig) *pb.ApplicationConfig {
	return &pb.ApplicationConfig{
		UseMariadb:  c.UseMariaDB,
		UseMongodb:  c.UseMongoDB,
		BuildType:   BuildTypeMapper.IntoMust(c.BuildType),
		BuildConfig: ToPBBuildConfig(c.BuildConfig),
	}
}
