package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBApplicationMetric(metric *domain.AppMetric) *pb.ApplicationMetric {
	return &pb.ApplicationMetric{
		Time:  timestamppb.New(metric.Time),
		Value: metric.Value,
	}
}
