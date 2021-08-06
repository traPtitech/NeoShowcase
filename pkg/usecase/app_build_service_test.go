package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	mock_pb "github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/mock"
	mock_repository "github.com/traPtitech/neoshowcase/pkg/interface/repository/mock"
)

func TestAppBuildService_QueueBuild(t *testing.T) {
	t.Parallel()
	t.Run("ビルドキューへの追加(Image)", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		repo := mock_repository.NewMockApplicationRepository(mockCtrl)
		c := mock_pb.NewMockBuilderServiceClient(mockCtrl)
		env := &domain.Environment{
			ID:            "5f34b184-9ae1-4969-95c0-0a016921d153",
			ApplicationID: "bee2466e-9d46-45e5-a6c4-4d359504c10c",
			BranchName:    "main",
			BuildType:     builder.BuildTypeImage,
		}
		res := &domain.Application{} //TODO: 正常時の値を取得する
		repo.EXPECT().GetApplicationByID(gomock.Any(), env.ApplicationID).Return(res, nil)
		s := NewAppBuildService(repo, c, "TestRegistry", "TestPrefix")
		err := s.QueueBuild(context.Background(), env)
		require.Nil(t, err)
	})

}
