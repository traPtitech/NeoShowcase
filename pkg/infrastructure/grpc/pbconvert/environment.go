package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBEnvironment(env *domain.Environment) *pb.ApplicationEnvVar {
	return &pb.ApplicationEnvVar{
		ApplicationId: env.ApplicationID,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}

func FromPBEnvironment(env *pb.ApplicationEnvVar) *domain.Environment {
	return &domain.Environment{
		ApplicationID: env.ApplicationId,
		Key:           env.Key,
		Value:         env.Value,
		System:        env.System,
	}
}
