package usecase

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type ContinuousDeploymentService interface {
	Run()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	bus      eventbus.Bus
	repo     repository.ApplicationRepository
	deployer AppDeployService
	builder  AppBuildService
}

func NewContinuousDeploymentService(bus eventbus.Bus, repo repository.ApplicationRepository, deployer AppDeployService, builder AppBuildService) ContinuousDeploymentService {
	return &continuousDeploymentService{
		bus:      bus,
		repo:     repo,
		deployer: deployer,
		builder:  builder,
	}
}

func (cd *continuousDeploymentService) Run() {
	cd.loop()
}

func (cd *continuousDeploymentService) Stop(ctx context.Context) error {
	return nil
}

func (cd *continuousDeploymentService) loop() {
	sub := cd.bus.Subscribe(event.BuilderBuildSucceeded, event.WebhookRepositoryPush)
	defer sub.Unsubscribe()
	for ev := range sub.Chan() {
		switch ev.Type {
		case event.WebhookRepositoryPush:
			repoURL := ev.Body["repository_url"].(string)
			branch := ev.Body["branch"].(string)
			cd.handleWebhookRepositoryPush(repoURL, branch)
		case event.BuilderBuildSucceeded:
			envID := ev.Body["environment_id"].(string)
			buildID := ev.Body["build_id"].(string)
			cd.handleBuilderBuildSucceeded(envID, buildID)
		}
	}
}

func (cd *continuousDeploymentService) handleWebhookRepositoryPush(repoURL string, branch string) {
	log.WithField("repo", repoURL).
		WithField("refs", branch).
		Info("repository push event received")

	env, err := cd.repo.GetEnvironmentByRepoAndBranch(context.Background(), repoURL, branch)
	if err != nil {
		if err == repository.ErrNotFound {
			return
		}
		log.WithError(err).
			WithField("repo", repoURL).
			WithField("refs", branch).
			Error("failed to GetEnvironmentByRepoAndBranch")
		return
	}

	err = cd.builder.QueueBuild(context.Background(), env)
	if err != nil {
		log.WithError(err).
			WithField("appID", env.ApplicationID).
			WithField("envID", env.ID).
			Error("failed to RequestBuild")
		return
	}
}

func (cd *continuousDeploymentService) handleBuilderBuildSucceeded(envID string, buildID string) {
	if envID == "" {
		// envIDが無い場合はテストビルド
		return
	}

	// 自動デプロイ
	log.WithField("envID", envID).
		WithField("buildID", buildID).
		Error("starting application")
	err := cd.deployer.QueueDeployment(context.Background(), envID, buildID)
	if err != nil {
		// TODO エラー処理
		log.WithError(err).
			WithField("envID", envID).
			WithField("buildID", buildID).
			Error("failed to Start Application")
	}
}
