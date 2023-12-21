package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var BuildStatusMapper = mapper.MustNewValueMapper(map[domain.BuildStatus]pb.BuildStatus{
	domain.BuildStatusQueued:    pb.BuildStatus_QUEUED,
	domain.BuildStatusBuilding:  pb.BuildStatus_BUILDING,
	domain.BuildStatusSucceeded: pb.BuildStatus_SUCCEEDED,
	domain.BuildStatusFailed:    pb.BuildStatus_FAILED,
	domain.BuildStatusCanceled:  pb.BuildStatus_CANCELLED,
	domain.BuildStatusSkipped:   pb.BuildStatus_SKIPPED,
})

func ToPBBuild(build *domain.Build) *pb.Build {
	return &pb.Build{
		Id:            build.ID,
		ApplicationId: build.ApplicationID,
		Commit:        build.Commit,
		Status:        BuildStatusMapper.IntoMust(build.Status),
		QueuedAt:      timestamppb.New(build.QueuedAt),
		StartedAt:     ToPBNullTimestamp(build.StartedAt),
		UpdatedAt:     ToPBNullTimestamp(build.UpdatedAt),
		FinishedAt:    ToPBNullTimestamp(build.FinishedAt),
		Retriable:     build.Retriable,
		Artifacts:     ds.Map(build.Artifacts, ToPBArtifact),
	}
}

func FromPBBuild(build *pb.Build) *domain.Build {
	return &domain.Build{
		ID:            build.Id,
		Commit:        build.Commit,
		ConfigHash:    "", /* Builder does not use this field */
		Status:        BuildStatusMapper.FromMust(build.Status),
		ApplicationID: build.ApplicationId,
		QueuedAt:      build.QueuedAt.AsTime(),
		StartedAt:     FromPBNullTimestamp(build.StartedAt),
		UpdatedAt:     FromPBNullTimestamp(build.UpdatedAt),
		FinishedAt:    FromPBNullTimestamp(build.FinishedAt),
		Retriable:     build.Retriable,
		Artifacts:     ds.Map(build.Artifacts, FromPBArtifact),
	}
}
