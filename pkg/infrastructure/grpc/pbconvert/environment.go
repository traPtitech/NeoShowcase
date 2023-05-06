package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBEnvironment(env *domain.Environment) *pb.ApplicationEnvVar {
	return &pb.ApplicationEnvVar{
		Key:    env.Key,
		Value:  env.Value,
		System: env.System,
	}
}
