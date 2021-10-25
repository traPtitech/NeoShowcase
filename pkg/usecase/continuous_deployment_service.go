package usecase

import (
	"context"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type ContinuousDeploymentService interface {
	Run()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	bus      domain.Bus
	repo     repository.ApplicationRepository
	deployer AppDeployService
	builder  AppBuildService
}

func NewContinuousDeploymentService(bus domain.Bus, repo repository.ApplicationRepository, deployer AppDeployService, builder AppBuildService) ContinuousDeploymentService {
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

	dbAdminUser := repoURL
	dbAdminPassword := generateRandomString(32)
	dbConfig := dbmanager.MariaDBConfig{
		"host",
		3307,
		dbAdminUser,
		dbAdminPassword,
	}

	db, err := dbmanager.NewMariaDBManager(dbConfig)
	if err != nil {
		log.WithError(err).
			WithField("Host", "host").
			WithField("Port", 3307).
			WithField("AdminUser", dbAdminUser).
			WithField("AdminPass", dbAdminPassword)
		return
	}

	dbSetting := domain.CreateArgs{
		Database: dbAdminUser,
		Password: dbAdminPassword,
	}
	err = db.Create(context.Background(), dbSetting)
	if err != nil {
		log.WithError(err).
			WithField("Database", dbSetting.Database).
			WithField("Password", dbSetting.Password)
	}

	branch, err := cd.repo.GetEnvironmentByRepoAndBranch(context.Background(), repoURL, branchName)
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

	_, err = cd.builder.QueueBuild(context.Background(), branch)
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
