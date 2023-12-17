package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBArtifact(artifact *domain.Artifact) *pb.Artifact {
	return &pb.Artifact{
		Id:        artifact.ID,
		Name:      artifact.Name,
		BuildId:   artifact.BuildID,
		Size:      artifact.Size,
		CreatedAt: timestamppb.New(artifact.CreatedAt),
		DeletedAt: ToPBNullTimestamp(artifact.DeletedAt),
	}
}

func FromPBArtifact(artifact *pb.Artifact) *domain.Artifact {
	return &domain.Artifact{
		ID:        artifact.Id,
		Name:      artifact.Name,
		BuildID:   artifact.BuildId,
		Size:      artifact.Size,
		CreatedAt: artifact.CreatedAt.AsTime(),
		DeletedAt: FromPBNullTimestamp(artifact.DeletedAt),
	}
}
