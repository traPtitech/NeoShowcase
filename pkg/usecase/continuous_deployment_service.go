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
	mariaDBManager domain.MariaDBManager
	mongoDBManager domain.MongoDBManager
}

func NewContinuousDeploymentService(bus domain.Bus, appRepo repository.ApplicationRepository, envRepo repository.EnvironmentRepository, deployer AppDeployService, builder AppBuildService, mariaDBManager domain.MariaDBManager, mongoDBManager domain.MongoDBManager) ContinuousDeploymentService {
	return &continuousDeploymentService{
		bus:            bus,
		appRepo:        appRepo,
		envRepo:        envRepo,
		deployer:       deployer,
		builder:        builder,
		mariaDBManager: mariaDBManager,
		mongoDBManager: mongoDBManager,
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
		case event.FetcherRequestApplicationBuild:
			applicationID := ev.Body["application_id"].(string)
			cd.handleNewBuildRequest(applicationID)
		case event.BuilderBuildSucceeded:
			branchID := ev.Body["application_id"].(string)
			buildID := ev.Body["build_id"].(string)
			cd.handleBuilderBuildSucceeded(branchID, buildID)
		}
	}
}

const (
	lowerCharSet  = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolCharSet = "!@#$%&*"
	numberSet     = "0123456789"
	allCharSet    = lowerCharSet + upperCharSet + symbolCharSet + numberSet
)

func generateRandomString(length int) string {
	var payload strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(allCharSet))
		payload.WriteByte(allCharSet[random])
	}
	return payload.String()
}

func (cd *continuousDeploymentService) handleNewBuildRequest(applicationID string) {
	log.WithField("applicationID", applicationID).
		Info("application build request event received")

	ctx := context.Background()

	app, err := cd.appRepo.GetApplicationByID(ctx, applicationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return
		}
		log.WithError(err).
			WithField("applicationID", applicationID).
			Error("failed to GetApplicationByID")
		return
	}

	// TODO dbSettingを設定から取得する
	dbName := fmt.Sprintf("%s_%s", app.Repository, app.ID)
	// TODO: アプリケーションの設定の取得
	applicationNeedsMariaDB := true
	if applicationNeedsMariaDB {
		dbExists, err := cd.mariaDBManager.IsExist(ctx, dbName)
		if err != nil {
			log.WithError(err).
				WithField("applicationID", applicationID).
				Error("failed to check if database exists")
			return
		}

		if !dbExists {
			dbPassword := generateRandomString(32)
			dbSetting := domain.CreateArgs{
				Database: dbName,
				Password: dbPassword,
			}

			if err := cd.mariaDBManager.Create(ctx, dbSetting); err != nil {
				log.WithError(err).
					WithField("Database", dbSetting.Database).
					WithField("Password", dbSetting.Password)
				return
			}

			if err := cd.envRepo.SetEnv(ctx, app.ID, domain.EnvMySQLUserKey, dbName); err != nil {
				log.WithError(err).
					WithField("applicationID", app.ID).
					WithField("Key", domain.EnvMySQLUserKey).
					WithField("Value", dbName)
				return
			}
			if err := cd.envRepo.SetEnv(ctx, app.ID, domain.EnvMySQLPasswordKey, dbPassword); err != nil {
				log.WithError(err).
					WithField("applicationID", app.ID).
					WithField("Key", domain.EnvMySQLPasswordKey)
				return
			}
			if err := cd.envRepo.SetEnv(ctx, app.ID, domain.EnvMySQLDatabaseKey, dbName); err != nil {
				log.WithError(err).
					WithField("applicationID", app.ID).
					WithField("Key", domain.EnvMySQLDatabaseKey).
					WithField("Value", dbName)
				return
			}
		}
	}

	// TODO: アプリケーションの設定の取得
	applicationNeedsMongoDB := true
	if applicationNeedsMongoDB {
		dbExists, err := cd.mongoDBManager.IsExist(ctx, dbName)
		if err != nil {
			log.WithError(err).
				WithField("applicationID", app.ID).
				WithField("dbName", dbName).
				Error("failed to check if database exists")
			return
		}

		if !dbExists {
			dbPassword := generateRandomString(32)
			dbSetting := domain.CreateArgs{
				Database: dbName,
				Password: dbPassword,
			}

			err := cd.mongoDBManager.Create(ctx, dbSetting)
			if err != nil {
				log.WithError(err).
					WithField("Database", dbSetting.Database).
					WithField("Password", dbSetting.Password)
			}

			if err := cd.envRepo.SetEnv(ctx, app.ID, domain.EnvMongoDBUserKey, dbName); err != nil {
				log.WithError(err).
					WithField("applicationID", app.ID).
					WithField("Key", domain.EnvMongoDBUserKey).
					WithField("Value", dbName)
				return
			}
			if err := cd.envRepo.SetEnv(ctx, app.ID, domain.EnvMongoDBPasswordKey, dbPassword); err != nil {
				log.WithError(err).
					WithField("applicationID", app.ID).
					WithField("Key", domain.EnvMongoDBPasswordKey)
				return
			}
			if err := cd.envRepo.SetEnv(ctx, app.ID, domain.EnvMongoDBDatabaseKey, dbName); err != nil {
				log.WithError(err).
					WithField("applicationID", app.ID).
					WithField("Key", domain.EnvMongoDBDatabaseKey).
					WithField("Value", dbName)
				return
			}
		}
	}

	_, err = cd.builder.QueueBuild(ctx, app)
	if err != nil {
		log.WithError(err).
			WithField("appID", app.ID).
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
