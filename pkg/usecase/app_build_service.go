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
	requestInterval = 10 * time.Second
)

type AppBuildService interface {
	QueueBuild(ctx context.Context, branch *domain.Branch) (string, error)
	CancelBuild(ctx context.Context, buildID string) error
	Shutdown()
}

type appBuildService struct {
	appRepo      repository.ApplicationRepository
	buildLogRepo repository.BuildLogRepository
	builder      pb.BuilderServiceClient

	queue           queue
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

func NewAppBuildService(appRepo repository.ApplicationRepository, buildLogRepo repository.BuildLogRepository, builder pb.BuilderServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString) AppBuildService {
	s := &appBuildService{
		appRepo:         appRepo,
		buildLogRepo:    buildLogRepo,
		builder:         builder,
		queue:           newQueue(),
		canceledJobList: []string{},
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
	go s.startQueueManager()
	return s
}

func (s *appBuildService) QueueBuild(ctx context.Context, branch *domain.Branch) (string, error) {
	app, err := s.appRepo.GetApplicationByID(ctx, branch.ApplicationID)
	if err != nil {
		return "", fmt.Errorf("Failed to QueueBuild: %w", err)
	}

	// ビルドログのエントリをDBに挿入
	buildLog, err := s.buildLogRepo.CreateBuildLog(ctx, branch.ID)
	if err != nil {
		log.WithError(err).Errorf("failed to create build log: %s", branch.ID)
		return "", fmt.Errorf("failed to create build log: %w", err)
	}

	s.queueWait.Add(1)
	s.queue.Push(&buildJob{
		buildID: buildLog.ID,
		app:     app,
		branch:  branch,
	})

	return buildLog.ID, nil
}

func (s *appBuildService) Shutdown() {
	s.queueWait.Wait()
}

func (s *appBuildService) CancelBuild(ctx context.Context, buildID string) error {
	if !s.isCanceled(buildID) {
		s.canceledJobList = append(s.canceledJobList, buildID)
		return nil
	}
	return fmt.Errorf("job is already canceled")
}

func (s *appBuildService) startQueueManager() {
	for {
		time.Sleep(requestInterval)

		v := s.queue.Top()
		if v == nil {
			continue
		}

		// キャンセルされたタスクならスキップ
		if s.isCanceled(v.buildID) {
			for i, id := range s.canceledJobList {
				if id == v.buildID {
					s.canceledJobList = append(s.canceledJobList[:i], s.canceledJobList[i+1:]...)
					s.queue.Pop()
					continue
				}
			}
		}

		res, err := s.builder.GetStatus(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.WithError(err).Error("failed to get status")
			break
		}
		if res.GetStatus() == pb.BuilderStatus_WAITING {
			s.queue.Pop()
			err := s.requestBuild(context.Background(), v.app, v.branch, v.buildID)
			if err != nil {
				log.WithError(err).Error("failed to request build")
			}
			s.queueWait.Done()
			break
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

type queue struct {
	data  []*buildJob
	mutex *sync.RWMutex
}

func newQueue() queue {
	return queue{
		data:  []*buildJob{},
		mutex: &sync.RWMutex{},
	}
}

func (q *queue) Push(job *buildJob) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.data = append(q.data, job)
}

func (q *queue) Top() *buildJob {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if len(q.data) == 0 {
		return nil
	}
	return q.data[0]
}

func (q *queue) Pop() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.data) == 0 {
		return
	}
	q.data = q.data[1:]
}
