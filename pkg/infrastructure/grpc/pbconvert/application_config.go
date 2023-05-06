package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func FromPBApplicationConfig(c *pb.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		UseMariaDB:  c.UseMariadb,
		UseMongoDB:  c.UseMongodb,
		BuildConfig: FromPBBuildConfig(c.BuildConfig),
		Entrypoint:  c.Entrypoint,
		Command:     c.Command,
	}
}

func ToPBApplicationConfig(c domain.ApplicationConfig) *pb.ApplicationConfig {
	return &pb.ApplicationConfig{
		UseMariadb:  c.UseMariaDB,
		UseMongodb:  c.UseMongoDB,
		BuildConfig: ToPBBuildConfig(c.BuildConfig),
		Entrypoint:  c.Entrypoint,
		Command:     c.Command,
	}
}
