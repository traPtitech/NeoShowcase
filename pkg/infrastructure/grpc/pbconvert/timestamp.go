package pbconvert

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func ToPBNullTimestamp(t optional.Of[time.Time]) *pb.NullTimestamp {
	return &pb.NullTimestamp{Timestamp: timestamppb.New(t.V), Valid: t.Valid}
}

func FromPBNullTimestamp(t *pb.NullTimestamp) optional.Of[time.Time] {
	if t.Valid {
		return optional.From(t.Timestamp.AsTime())
	}
	return optional.None[time.Time]()
}
