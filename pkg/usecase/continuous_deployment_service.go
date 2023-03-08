package usecase

import (
	"context"
	"errors"
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
	bus      domain.Bus
	appRepo  repository.ApplicationRepository
	envRepo  repository.EnvironmentRepository
	deployer AppDeployService
	builder  AppBuildService
}

func NewContinuousDeploymentService(bus domain.Bus, appRepo repository.ApplicationRepository, envRepo repository.EnvironmentRepository, deployer AppDeployService, builder AppBuildService) ContinuousDeploymentService {
	return &continuousDeploymentService{
		bus:      bus,
		appRepo:  appRepo,
		envRepo:  envRepo,
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
