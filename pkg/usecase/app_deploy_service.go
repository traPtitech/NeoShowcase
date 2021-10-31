package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AppDeployService interface {
	QueueDeployment(ctx context.Context, envID string, buildID string) error
}

type appDeployService struct {
	backend domain.Backend
	ss      pb.StaticSiteServiceClient

	imageRegistry   string
	imageNamePrefix string

	// TODO 後で消す
	db *sql.DB
}

func NewAppDeployService(backend domain.Backend, ss pb.StaticSiteServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString, db *sql.DB) AppDeployService {
	return &appDeployService{
		backend:         backend,
		ss:              ss,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
		db:              db,
	}
}

func (s *appDeployService) QueueDeployment(ctx context.Context, branchID string, buildID string) error {
	branch, err := models.Branches(
		qm.Load(models.BranchRels.Website),
		models.BranchWhere.ID.EQ(branchID),
	).One(ctx, s.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("branch (%s) was not found", branchID)
		}
		return fmt.Errorf("failed to get Branch: %w", err)
	}

	ok, err := branch.BuildLogs(
		models.BuildLogWhere.ID.EQ(buildID),
		models.BuildLogWhere.Result.EQ(models.BuildLogsResultSUCCEEDED),
	).Exists(ctx, s.db)
	if err != nil {
		return fmt.Errorf("failed to BuildLogExists: %w", err)
	}
	if !ok {
		return fmt.Errorf("build (%s) was not found", buildID)
	}

	entry := &appDeployment{
		BranchID: branchID,
		BuildID:  buildID,
		QueuedAt: time.Now(),
		branch:   branch,
	}
	// TODO ちゃんとキューに入れて非同期処理する
	go s.deploy(entry)
	return nil
}

type appDeployment struct {
	BranchID string
	BuildID  string
	QueuedAt time.Time
	branch   *models.Branch
}

func (s *appDeployService) deploy(entry *appDeployment) {
	branch := entry.branch
	website := branch.R.Website

	// TODO Ingressの設定がここでする必要があるかどうか確認する
	switch builder.BuildTypeFromString(branch.BuildType) {
	case builder.BuildTypeImage:
		var httpProxy *domain.ContainerHTTPProxy
		if website != nil {
			httpProxy = &domain.ContainerHTTPProxy{
				Domain: website.FQDN,
				Port:   website.HTTPPort,
			}
		}

		err := s.backend.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ApplicationID: branch.ApplicationID,
			EnvironmentID: branch.ID,
			ImageName:     builder.GetImageName(s.imageRegistry, s.imageNamePrefix, branch.ApplicationID),
			ImageTag:      entry.BuildID,
			HTTPProxy:     httpProxy,
			Recreate:      true,
		})
		if err != nil {
			log.WithField("BranchID", entry.BranchID).
				WithField("buildID", entry.BuildID).
				WithField("queuedAt", entry.QueuedAt).
				WithError(err).
				Errorf("failed to create container")
			return
		}

		branch.BuildID = null.StringFrom(entry.BuildID)
		if _, err := branch.Update(context.Background(), s.db, boil.Infer()); err != nil {
			log.WithField("branchID", entry.BranchID).
				WithField("buildID", entry.BuildID).
				WithField("queuedAt", entry.QueuedAt).
				WithError(err).
				Errorf("failed to update branch")
			return
		}

	case builder.BuildTypeStatic:
		branch.BuildID = null.StringFrom(entry.BuildID)
		if _, err := branch.Update(context.Background(), s.db, boil.Infer()); err != nil {
			log.WithField("branchID", entry.BranchID).
				WithField("buildID", entry.BuildID).
				WithField("queuedAt", entry.QueuedAt).
				WithError(err).
				Errorf("failed to update env")
			return
		}

		if _, err := s.ss.Reload(context.Background(), &pb.ReloadRequest{}); err != nil {
			log.WithField("branchID", entry.BranchID).
				WithField("buildID", entry.BuildID).
				WithField("queuedAt", entry.QueuedAt).
				WithError(err).
				Errorf("failed to reload StaticSiteServer")
			return
		}
	}
}
