package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPBRuntimeImage(image *domain.RuntimeImage) *pb.RuntimeImage {
	return &pb.RuntimeImage{
		BuildId:   image.BuildID,
		Size:      image.Size,
		CreatedAt: timestamppb.New(image.CreatedAt),
	}
}
