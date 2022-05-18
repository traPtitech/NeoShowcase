package usecase

import (
	"context"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type ContinuousDeploymentService interface {
	Run()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	bus            domain.Bus
	repo           repository.ApplicationRepository
	deployer       AppDeployService
	builder        AppBuildService
	mariadbmanager domain.MariaDBManager
	mongodbmanager domain.MongoManager
}

func NewContinuousDeploymentService(bus domain.Bus, repo repository.ApplicationRepository, deployer AppDeployService, builder AppBuildService, mariadbmanager domain.MariaDBManager, mongodbmanager domain.MongoManager) ContinuousDeploymentService {
	return &continuousDeploymentService{
		bus:            bus,
		repo:           repo,
		deployer:       deployer,
		builder:        builder,
		mariadbmanager: mariadbmanager,
		mongodbmanager: mongodbmanager,
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
			branchID := ev.Body["branch_id"].(string)
			buildID := ev.Body["build_id"].(string)
			cd.handleBuilderBuildSucceeded(branchID, buildID)
		}
	}
}

func generateRandomString(length int) string {
	lowerCharSet := "abcdedfghijklmnopqrst"
	upperCharSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolCharSet := "!@#$%&*"
	numberSet := "0123456789"
	allCharSet := lowerCharSet + upperCharSet + symbolCharSet + numberSet

	var payload strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(allCharSet))
		payload.WriteByte(allCharSet[random])
	}

	return payload.String()
}

func (cd *continuousDeploymentService) handleWebhookRepositoryPush(repoURL string, branchName string) {
	log.WithField("repo", repoURL).
		WithField("refs", branchName).
		Info("repository push event received")

	ctx := context.Background()

	dbName := repoURL + branchName
	// TODO: アプリケーションの設定の取得
	applicationNeedsMariaDB := true
	if applicationNeedsMariaDB {
		dbExists, err := cd.mariadbmanager.IsExist(context.Background(), dbName)
		if err != nil {
			log.WithError(err).
				WithField("repo", repoURL).
				WithField("refs", branchName).
				Error("failed to check if database exists")
			return
		}

		if !dbExists {
			// TODO dbUser, dbSettingを設定から取得する
			dbUser := repoURL
			dbPassword := generateRandomString(32)
			dbSetting := domain.CreateArgs{
				Database: dbUser,
				Password: dbPassword,
			}

			err := cd.mariadbmanager.Create(ctx, dbSetting)
			if err != nil {
				log.WithError(err).
					WithField("Database", dbSetting.Database).
					WithField("Password", dbSetting.Password)
			}
		}
	}

	// TODO: アプリケーションの設定の取得
	applicationNeedsMongoDB := true
	if applicationNeedsMongoDB {
		dbExists, err := cd.mongodbmanager.IsExist(context.Background(), dbName)
		if err != nil {
			log.WithError(err).
				WithField("repo", repoURL).
				WithField("refs", branchName).
				Error("failed to check if database exists")
			return
		}

		if !dbExists {
			// TODO dbUser, dbSettingを設定から取得する
			dbUser := repoURL
			dbPassword := generateRandomString(32)
			dbSetting := domain.CreateArgs{
				Database: dbUser,
				Password: dbPassword,
			}

			err := cd.mongodbmanager.Create(ctx, dbSetting)
			if err != nil {
				log.WithError(err).
					WithField("Database", dbSetting.Database).
					WithField("Password", dbSetting.Password)
			}
		}
	}

	branch, err := cd.repo.GetBranchByRepoAndBranchName(ctx, repoURL, branchName)
	if err != nil {
		if err == repository.ErrNotFound {
			return
		}
		log.WithError(err).
			WithField("repo", repoURL).
			WithField("refs", branchName).
			Error("failed to GetBranchByRepoAndBranchName")
		return
	}

	_, err = cd.builder.QueueBuild(ctx, branch)
	if err != nil {
		log.WithError(err).
			WithField("appID", branch.ApplicationID).
			WithField("branchID", branch.ID).
			Error("failed to RequestBuild")
		return
	}
}

func (cd *continuousDeploymentService) handleBuilderBuildSucceeded(branchID string, buildID string) {
	if branchID == "" {
		// branchIDが無い場合はテストビルド
		return
	}

	// 自動デプロイ
	log.WithField("branchID", branchID).
		WithField("buildID", buildID).
		Error("starting application")
	err := cd.deployer.QueueDeployment(context.Background(), branchID, buildID)
	if err != nil {
		// TODO エラー処理
		log.WithError(err).
			WithField("branchID", branchID).
			WithField("buildID", buildID).
			Error("failed to Start Application")
	}
}
