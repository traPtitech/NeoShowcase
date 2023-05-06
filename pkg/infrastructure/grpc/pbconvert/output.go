package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBApplicationOutput(l *domain.ContainerLog) *pb.ApplicationOutput {
	return &pb.ApplicationOutput{
		Time: timestamppb.New(l.Time),
		Log:  l.Log,
	}
}
