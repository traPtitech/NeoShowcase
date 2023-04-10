package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ToPBNullTimestamp(t optional.Of[time.Time]) *pb.NullTimestamp {
	return &pb.NullTimestamp{Timestamp: timestamppb.New(t.V), Valid: t.Valid}
}
