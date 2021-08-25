package repository

import (
	"context"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type ApplicationRepository interface {
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	GetApplicationByID(ctx context.Context, id string) (*domain.Application, error)
	CreateEnvironment(ctx context.Context, appID string, branchName string, buildType builder.BuildType) (*domain.Environment, error)
	GetEnvironmentByID(ctx context.Context, id string) (*domain.Environment, error)
	GetEnvironmentByRepoAndBranch(ctx context.Context, repoURL string, branch string) (*domain.Environment, error)
	SetWebsite(ctx context.Context, envID string, fqdn string, httpPort int) error
}

type applicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

type CreateApplicationArgs struct {
	Owner         string
	Name          string
	RepositoryURL string
	BranchName    string
	BuildType     builder.BuildType
}

func (r *applicationRepository) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	const errMsg = "failed to CreateApplication: %w"

	// リポジトリ情報を設定
	repo, err := models.Repositories(models.RepositoryWhere.Remote.EQ(args.RepositoryURL)).One(ctx, r.db)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf(errMsg, err)
	} else if repo == nil {
		repo = &models.Repository{
			ID:     domain.NewID(),
			Remote: args.RepositoryURL,
		}
		if err := repo.Insert(ctx, r.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}

	// アプリケーション作成
	app := &models.Application{
		ID:           domain.NewID(),
		Owner:        args.Owner,
		Name:         args.Name,
		RepositoryID: repo.ID,
	}
	if err := app.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("appID", app.ID).
		Info("app created")

	// 初期Env作成
	env := &models.Environment{
		ID:         domain.NewID(),
		BranchName: args.BranchName,
		BuildType:  args.BuildType.String(),
	}
	if err := app.AddEnvironments(ctx, r.db, true, env); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("appID", app.ID).
		WithField("envID", env.ID).
		Info("env created")

	return &domain.Application{
		ID: app.ID,
		Repository: domain.Repository{
			ID:        repo.ID,
			RemoteURL: repo.Remote,
		},
	}, nil
}

func (r *applicationRepository) GetApplicationByID(ctx context.Context, id string) (*domain.Application, error) {
	const errMsg = "failed to GetApplicationByID: %w"

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		models.ApplicationWhere.ID.EQ(id),
		models.ApplicationWhere.DeletedAt.IsNull(),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	return &domain.Application{
		ID: app.ID,
		Repository: domain.Repository{
			ID:        app.R.Repository.ID,
			RemoteURL: app.R.Repository.Remote,
		},
	}, nil
}

func (r *applicationRepository) CreateEnvironment(ctx context.Context, appID string, branchName string, buildType builder.BuildType) (*domain.Environment, error) {
	const errMsg = "failed to CreateEnvironment: %w"

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Environments, models.EnvironmentWhere.BranchName.EQ(branchName)),
		models.ApplicationWhere.ID.EQ(appID),
		models.ApplicationWhere.DeletedAt.IsNull(),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	// 指定したブランチのEnvironmentが存在しないことを確認
	if len(app.R.Environments) > 0 {
		return nil, fmt.Errorf("the environment for branch `%s` has already existed", branchName)
	}

	env := &models.Environment{
		ID:         domain.NewID(),
		BranchName: branchName,
		BuildType:  buildType.String(),
	}
	if err := app.AddEnvironments(ctx, r.db, true, env); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("appID", env.ApplicationID).
		WithField("envID", env.ID).
		Info("env created")

	return &domain.Environment{
		ID:            env.ID,
		ApplicationID: env.ApplicationID,
		BranchName:    env.BranchName,
		BuildType:     builder.BuildTypeFromString(env.BuildType),
	}, nil
}

func (r *applicationRepository) GetEnvironmentByID(ctx context.Context, id string) (*domain.Environment, error) {
	const errMsg = "failed to GetEnvironmentByID: %w"

	env, err := models.FindEnvironment(ctx, r.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	return &domain.Environment{
		ID:            env.ID,
		ApplicationID: env.ApplicationID,
		BranchName:    env.BranchName,
		BuildType:     builder.BuildTypeFromString(env.BuildType),
	}, nil
}

func (r *applicationRepository) GetEnvironmentByRepoAndBranch(ctx context.Context, repoURL string, branch string) (*domain.Environment, error) {
	const errMsg = "failed to GetEnvironmentByRepoAndBranch: %w"

	repo, err := models.Repositories(
		models.RepositoryWhere.Remote.EQ(repoURL),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	app, err := repo.Applications(
		qm.Load(models.ApplicationRels.Environments, models.EnvironmentWhere.BranchName.EQ(branch)),
		models.ApplicationWhere.DeletedAt.IsNull(),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	if len(app.R.Environments) > 0 {
		env := app.R.Environments[0]
		return &domain.Environment{
			ID:            env.ID,
			ApplicationID: env.ApplicationID,
			BranchName:    env.BranchName,
			BuildType:     builder.BuildTypeFromString(env.BuildType),
		}, nil
	}
	return nil, ErrNotFound
}

func (r *applicationRepository) SetWebsite(ctx context.Context, envID string, fqdn string, httpPort int) error {
	const errMsg = "failed to SetWebsite: %w"

	env, err := models.Environments(
		qm.Load(models.EnvironmentRels.Website),
		models.EnvironmentWhere.ID.EQ(envID),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf(errMsg, err)
	}

	ws := env.R.Website
	if ws != nil {
		// テーブルの情報を更新
		ws.FQDN = fqdn
		ws.HTTPPort = httpPort
		if _, err := ws.Update(ctx, r.db, boil.Infer()); err != nil {
			return fmt.Errorf(errMsg, err)
		}
		return nil
	}

	// Websiteをテーブルに挿入
	ws = &models.Website{
		ID:       domain.NewID(),
		FQDN:     fqdn,
		HTTPPort: httpPort,
	}
	if err := env.SetWebsite(ctx, r.db, true, ws); err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}
