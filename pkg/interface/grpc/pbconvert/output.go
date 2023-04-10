package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPBApplicationOutput(l *domain.ContainerLog) *pb.ApplicationOutput {
	return &pb.ApplicationOutput{
		Time: timestamppb.New(l.Time),
		Log:  l.Log,
	}
}
