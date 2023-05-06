package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBArtifact(artifact *domain.Artifact) *pb.Artifact {
	return &pb.Artifact{
		Id:        artifact.ID,
		Size:      artifact.Size,
		CreatedAt: timestamppb.New(artifact.CreatedAt),
		DeletedAt: ToPBNullTimestamp(artifact.DeletedAt),
	}
}
