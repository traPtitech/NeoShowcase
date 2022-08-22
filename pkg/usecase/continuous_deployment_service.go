package usecase

import (
	"context"
	"errors"
	"fmt"
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
	appRepo        repository.ApplicationRepository
	envRepo        repository.EnvironmentRepository
	deployer       AppDeployService
	builder        AppBuildService
	mariadbmanager domain.MariaDBManager
	mongodbmanager domain.MongoDBManager
}

func NewContinuousDeploymentService(bus domain.Bus, appRepo repository.ApplicationRepository, envRepo repository.EnvironmentRepository, deployer AppDeployService, builder AppBuildService, mariadbmanager domain.MariaDBManager, mongodbmanager domain.MongoDBManager) ContinuousDeploymentService {
	return &continuousDeploymentService{
		bus:            bus,
		appRepo:        appRepo,
		envRepo:        envRepo,
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

	branch, err := cd.appRepo.GetBranchByRepoAndBranchName(ctx, repoURL, branchName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return
		}
		log.WithError(err).
			WithField("repo", repoURL).
			WithField("refs", branchName).
			Error("failed to GetBranchByRepoAndBranchName")
		return
	}

	// TODO dbSettingを設定から取得する
	dbName := fmt.Sprintf("%s_%s", repoURL, branchName)
	// TODO: アプリケーションの設定の取得
	applicationNeedsMariaDB := true
	if applicationNeedsMariaDB {
		dbExists, err := cd.mariadbmanager.IsExist(ctx, dbName)
		if err != nil {
			log.WithError(err).
				WithField("repo", repoURL).
				WithField("refs", branchName).
				Error("failed to check if database exists")
			return
		}

		if !dbExists {
			dbPassword := generateRandomString(32)
			dbSetting := domain.CreateArgs{
				Database: dbName,
				Password: dbPassword,
			}

			if err := cd.mariadbmanager.Create(ctx, dbSetting); err != nil {
				log.WithError(err).
					WithField("Database", dbSetting.Database).
					WithField("Password", dbSetting.Password)
				return
			}

			if err := cd.envRepo.SetEnv(ctx, branch.ID, domain.EnvMySQLDatabaseKey, dbName); err != nil {
				log.WithError(err).
					WithField("BranchID", branch.ID).
					WithField("BranchName", branchName).
					WithField("Key", domain.EnvMySQLDatabaseKey)
				return
			}
			if err := cd.envRepo.SetEnv(ctx, branch.ID, domain.EnvMySQLPasswordKey, dbPassword); err != nil {
				log.WithError(err).
					WithField("BranchID", branch.ID).
					WithField("BranchName", branchName).
					WithField("Key", domain.EnvMySQLPasswordKey)
				return
			}
		}
	}

	// TODO: アプリケーションの設定の取得
	applicationNeedsMongoDB := true
	if applicationNeedsMongoDB {
		dbExists, err := cd.mongodbmanager.IsExist(ctx, dbName)
		if err != nil {
			log.WithError(err).
				WithField("repo", repoURL).
				WithField("refs", branchName).
				Error("failed to check if database exists")
			return
		}

		if !dbExists {
			dbPassword := generateRandomString(32)
			dbSetting := domain.CreateArgs{
				Database: dbName,
				Password: dbPassword,
			}

			err := cd.mongodbmanager.Create(ctx, dbSetting)
			if err != nil {
				log.WithError(err).
					WithField("Database", dbSetting.Database).
					WithField("Password", dbSetting.Password)
			}

			if err := cd.envRepo.SetEnv(ctx, branch.ID, domain.EnvMongoDBDatabaseKey, dbName); err != nil {
				log.WithError(err).
					WithField("BranchID", branch.ID).
					WithField("BranchName", branchName).
					WithField("Key", domain.EnvMongoDBDatabaseKey)
				return
			}
			if err := cd.envRepo.SetEnv(ctx, branch.ID, domain.EnvMongoDBPasswordKey, dbPassword); err != nil {
				log.WithError(err).
					WithField("BranchID", branch.ID).
					WithField("BranchName", branchName).
					WithField("Key", domain.EnvMongoDBPasswordKey)
				return
			}
		}
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
