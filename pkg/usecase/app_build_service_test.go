package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	mock_pb "github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/mock"
	mock_repository "github.com/traPtitech/neoshowcase/pkg/interface/repository/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAppBuildService_QueueBuild(t *testing.T) {
	t.Run("ビルドキューへの追加(Image)", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		repo := mock_repository.NewMockApplicationRepository(mockCtrl)
		c := mock_pb.NewMockBuilderServiceClient(mockCtrl)
		s := NewAppBuildService(repo, c, "TestRegistry", "TestPrefix")

		env := &domain.Environment{
			ID:            "5f34b184-9ae1-4969-95c0-0a016921d153",
			ApplicationID: "bee2466e-9d46-45e5-a6c4-4d359504c10c",
			BranchName:    "main",
			BuildType:     builder.BuildTypeImage,
		}
		res := &domain.Application{
			Repository: domain.Repository{
				RemoteURL: "https://git.trap.jp/hijiki51/git-test",
			},
		}

		repo.EXPECT().
			GetApplicationByID(context.Background(), env.ApplicationID).Return(res, nil)

		c.EXPECT().
			GetStatus(context.Background(), &emptypb.Empty{}).
			Return(&pb.GetStatusResponse{Status: pb.BuilderStatus_WAITING}, nil).
			AnyTimes()

		c.EXPECT().
			StartBuildImage(context.Background(), &pb.StartBuildImageRequest{
				ImageName: "TestRegistry/TestPrefixbee2466e-9d46-45e5-a6c4-4d359504c10c",
				Source: &pb.BuildSource{
					RepositoryUrl: res.Repository.RemoteURL,
				},
				Options:       &pb.BuildOptions{},
				EnvironmentId: env.ID,
			}).
			Return(&pb.StartBuildImageResponse{}, nil)

		err := s.QueueBuild(context.Background(), env)
		s.Shutdown()
		require.Nil(t, err)
	})
}
