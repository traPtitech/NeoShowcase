package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
)

var BuildStatusMapper = mapper.MustNewValueMapper(map[domain.BuildStatus]pb.Build_BuildStatus{
	domain.BuildStatusQueued:    pb.Build_QUEUED,
	domain.BuildStatusBuilding:  pb.Build_BUILDING,
	domain.BuildStatusSucceeded: pb.Build_SUCCEEDED,
	domain.BuildStatusFailed:    pb.Build_FAILED,
	domain.BuildStatusCanceled:  pb.Build_CANCELLED,
	domain.BuildStatusSkipped:   pb.Build_SKIPPED,
})

func ToPBBuild(build *domain.Build) *pb.Build {
	b := &pb.Build{
		Id:         build.ID,
		Commit:     build.Commit,
		Status:     BuildStatusMapper.IntoMust(build.Status),
		StartedAt:  ToPBNullTimestamp(build.StartedAt),
		UpdatedAt:  ToPBNullTimestamp(build.UpdatedAt),
		FinishedAt: ToPBNullTimestamp(build.FinishedAt),
		Retriable:  build.Retriable,
	}
	if build.Artifact.Valid {
		b.Artifact = ToPBArtifact(&build.Artifact.V)
	}
	return b
}
