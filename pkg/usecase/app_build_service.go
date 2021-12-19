package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	queueBufferSize = 10
	requestInterval = 10 * time.Second
)

type AppBuildService interface {
	QueueBuild(ctx context.Context, branch *domain.Branch) (string, error)
	CancelBuild(ctx context.Context, buildID string) error
	Shutdown()
}

type appBuildService struct {
	repo    repository.ApplicationRepository
	builder pb.BuilderServiceClient

	queue           chan *buildJob
	queueWait       sync.WaitGroup
	canceledJobList []string
	imageRegistry   string
	imageNamePrefix string
}

type buildJob struct {
	buildID string
	app     *domain.Application
	branch  *domain.Branch
}

var (
	ErrQueueFull = fmt.Errorf("queue is full")
)

func NewAppBuildService(repo repository.ApplicationRepository, builder pb.BuilderServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString) AppBuildService {
	s := &appBuildService{
		repo:            repo,
		builder:         builder,
		queue:           make(chan *buildJob, queueBufferSize),
		canceledJobList: []string{},
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
	go s.startQueueManager()
	return s
}

func (s *appBuildService) QueueBuild(ctx context.Context, branch *domain.Branch) (string, error) {
	app, err := s.repo.GetApplicationByID(ctx, branch.ApplicationID)
	if err != nil {
		return "", fmt.Errorf("Failed to QueueBuild: %w", err)
	}

	buildID := domain.NewID()
	// TODO: このタイミングでDBに入れる

	s.queueWait.Add(1)
	select {
	case s.queue <- &buildJob{
		buildID: buildID,
		app:     app,
		branch:  branch,
	}:
	default:
		return "", ErrQueueFull
	}

	return buildID, nil
}

func (s *appBuildService) Shutdown() {
	s.queueWait.Wait()
	// TODO: shutdown後にキューにジョブが追加されないことを保証しないといけない(でないとclosed channelにpushしようとしてpanicする)
	close(s.queue)
}

func (s *appBuildService) CancelBuild(ctx context.Context, buildID string) error {
	if !s.isCanceled(buildID) {
		s.canceledJobList = append(s.canceledJobList, buildID)
		return nil
	}
	return fmt.Errorf("job is already canceled")
}

func (s *appBuildService) startQueueManager() {
	for v := range s.queue {
		// キャンセルされたタスクならスキップ
		if s.isCanceled(v.buildID) {
			for i, id := range s.canceledJobList {
				if id == v.buildID {
					s.canceledJobList = append(s.canceledJobList[:i], s.canceledJobList[i+1:]...)
					continue
				}
			}
		}
		for {
			res, err := s.builder.GetStatus(context.Background(), &emptypb.Empty{})
			if err != nil {
				log.WithError(err).Error("failed to get status")
				break
			}
			if res.GetStatus() == pb.BuilderStatus_WAITING {
				err := s.requestBuild(context.Background(), v.app, v.branch, v.buildID)
				if err != nil {
					log.WithError(err).Error("failed to request build")
				}
				s.queueWait.Done()
				break
			}
			time.Sleep(requestInterval)
		}
	}
}

func (s *appBuildService) requestBuild(ctx context.Context, app *domain.Application, branch *domain.Branch, buildID string) error {
	switch branch.BuildType {
	case builder.BuildTypeImage:
		_, err := s.builder.StartBuildImage(ctx, &pb.StartBuildImageRequest{
			ImageName: builder.GetImageName(s.imageRegistry, s.imageNamePrefix, branch.ApplicationID),
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.RemoteURL, // TODO ブランチ・タグ指定に対応
			},
			Options:  &pb.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			BranchId: branch.ID,
			BuildId:  buildID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build image: %w", err)
		}

	case builder.BuildTypeStatic:
		_, err := s.builder.StartBuildStatic(ctx, &pb.StartBuildStaticRequest{
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.RemoteURL, // TODO ブランチ・タグ指定に対応
			},
			Options:  &pb.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			BranchId: branch.ID,
			BuildId:  buildID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build static: %w", err)
		}

	default:
		return fmt.Errorf("unknown build type: %s", branch.BuildType)
	}

	log.WithField("branchID", branch.ID).
		Info("build requested")
	return nil
}

func (s *appBuildService) isCanceled(buildID string) bool {
	for _, v := range s.canceledJobList {
		if v == buildID {
			return true
		}
	}
	return false
}
