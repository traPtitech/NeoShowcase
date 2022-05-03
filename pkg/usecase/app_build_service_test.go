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
	t.Parallel()

	t.Run("ビルドキューへの追加(Image)", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		appRepo := mock_repository.NewMockApplicationRepository(mockCtrl)
		buildLogRepo := mock_repository.NewMockBuildLogRepository(mockCtrl)
		c := mock_pb.NewMockBuilderServiceClient(mockCtrl)
		s := NewAppBuildService(appRepo, buildLogRepo, c, "TestRegistry", "TestPrefix")
		branch := &domain.Branch{
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
		buildLog := &domain.BuildLog{
			ID:       "f01691dd-985a-48c9-8b47-205af468431a",
			Result:   builder.BuildStatusQueued,
			BranchID: branch.ID,
		}

		appRepo.EXPECT().
			GetApplicationByID(context.Background(), branch.ApplicationID).Return(res, nil)

		buildLogRepo.EXPECT().CreateBuildLog(context.Background(), branch.ID).Return(buildLog, nil)

		c.EXPECT().
			GetStatus(context.Background(), &emptypb.Empty{}).
			Return(&pb.GetStatusResponse{
				Status:  pb.BuilderStatus_WAITING,
				BuildId: buildLog.ID,
			}, nil).
			AnyTimes()

		c.EXPECT().
			StartBuildImage(context.Background(), &pb.StartBuildImageRequest{
				ImageName: "TestRegistry/TestPrefixbee2466e-9d46-45e5-a6c4-4d359504c10c",
				Source: &pb.BuildSource{
					RepositoryUrl: res.Repository.RemoteURL,
				},
				Options:  &pb.BuildOptions{},
				BuildId:  buildLog.ID,
				BranchId: branch.ID,
			}).
			Return(&pb.StartBuildImageResponse{}, nil)

		_, err := s.QueueBuild(context.Background(), branch)
		s.Shutdown()
		require.Nil(t, err)
	})

	t.Run("ビルドキューへの追加(Static)", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		appRepo := mock_repository.NewMockApplicationRepository(mockCtrl)
		buildLogRepo := mock_repository.NewMockBuildLogRepository(mockCtrl)
		c := mock_pb.NewMockBuilderServiceClient(mockCtrl)
		s := NewAppBuildService(appRepo, buildLogRepo, c, "TestRegistry", "TestPrefix")
		branch := &domain.Branch{
			ID:            "1d9cc06d-813f-4cf7-947e-546e1a814fed",
			ApplicationID: "d563e2de-7905-4267-8a9c-51520aac02b3",
			BranchName:    "develop",
			BuildType:     builder.BuildTypeStatic,
		}
		res := &domain.Application{
			Repository: domain.Repository{
				RemoteURL: "https://git.trap.jp/hijiki51/git-test",
			},
		}
		buildLog := &domain.BuildLog{
			ID:       "f01691dd-985a-48c9-8b47-205af468431a",
			Result:   builder.BuildStatusQueued,
			BranchID: branch.ID,
		}

		appRepo.EXPECT().
			GetApplicationByID(context.Background(), branch.ApplicationID).Return(res, nil)

		buildLogRepo.EXPECT().CreateBuildLog(context.Background(), branch.ID).Return(buildLog, nil)

		c.EXPECT().
			GetStatus(context.Background(), &emptypb.Empty{}).
			Return(&pb.GetStatusResponse{
				Status:  pb.BuilderStatus_WAITING,
				BuildId: buildLog.ID,
			}, nil).
			AnyTimes()

		c.EXPECT().
			StartBuildStatic(context.Background(), &pb.StartBuildStaticRequest{
				Source: &pb.BuildSource{
					RepositoryUrl: res.Repository.RemoteURL,
				},
				Options:  &pb.BuildOptions{},
				BuildId:  buildLog.ID,
				BranchId: branch.ID,
			}).
			Return(&pb.StartBuildStaticResponse{}, nil)

		_, err := s.QueueBuild(context.Background(), branch)
		s.Shutdown()
		require.Nil(t, err)
	})
	t.Run("追加されたジョブのキャンセル", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		appRepo := mock_repository.NewMockApplicationRepository(mockCtrl)
		buildLogRepo := mock_repository.NewMockBuildLogRepository(mockCtrl)
		c := mock_pb.NewMockBuilderServiceClient(mockCtrl)
		s := NewAppBuildService(appRepo, buildLogRepo, c, "TestRegistry", "TestPrefix")
		branch := &domain.Branch{
			ID:            "1d9cc06d-813f-4cf7-947e-546e1a814fed",
			ApplicationID: "d563e2de-7905-4267-8a9c-51520aac02b3",
			BranchName:    "develop",
			BuildType:     builder.BuildTypeStatic,
		}
		res := &domain.Application{
			Repository: domain.Repository{
				RemoteURL: "https://git.trap.jp/hijiki51/git-test",
			},
		}
		buildLog := &domain.BuildLog{
			ID:       "f01691dd-985a-48c9-8b47-205af468431a",
			Result:   builder.BuildStatusQueued,
			BranchID: branch.ID,
		}

		appRepo.EXPECT().
			GetApplicationByID(context.Background(), branch.ApplicationID).Return(res, nil)

		buildLogRepo.EXPECT().CreateBuildLog(context.Background(), branch.ID).Return(buildLog, nil)

		c.EXPECT().
			GetStatus(context.Background(), &emptypb.Empty{}).
			Return(&pb.GetStatusResponse{
				Status:  pb.BuilderStatus_WAITING,
				BuildId: buildLog.ID,
			}, nil).
			AnyTimes()

		id, err := s.QueueBuild(context.Background(), branch)
		if err != nil {
			t.Fatal(err)
		}
		err = s.CancelBuild(context.Background(), id)
		s.Shutdown()
		require.Nil(t, err)
	})
}
