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
	CreateBranch(ctx context.Context, appID string, branchName string, buildType builder.BuildType) (*domain.Branch, error)
	GetBranchByID(ctx context.Context, id string) (*domain.Branch, error)
	GetBranchByRepoAndBranchName(ctx context.Context, repoURL string, branch string) (*domain.Branch, error)
	SetWebsite(ctx context.Context, branchID string, fqdn string, httpPort int) error
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

	// 初期Branch作成
	branch := &models.Branch{
		ID:         domain.NewID(),
		BranchName: args.BranchName,
		BuildType:  args.BuildType.String(),
	}
	if err := app.AddBranches(ctx, r.db, true, branch); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("appID", app.ID).
		WithField("branchID", branch.ID).
		Info("branch created")

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

func (r *applicationRepository) CreateBranch(ctx context.Context, appID string, branchName string, buildType builder.BuildType) (*domain.Branch, error) {
	const errMsg = "failed to CreateBranch: %w"

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Branches, models.BranchWhere.BranchName.EQ(branchName)),
		models.ApplicationWhere.ID.EQ(appID),
		models.ApplicationWhere.DeletedAt.IsNull(),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	// 指定したブランチが存在しないことを確認
	if len(app.R.Branches) > 0 {
		return nil, fmt.Errorf("the branch `%s` has already existed", branchName)
	}

	branch := &models.Branch{
		ID:         domain.NewID(),
		BranchName: branchName,
		BuildType:  buildType.String(),
	}
	if err := app.AddBranches(ctx, r.db, true, branch); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("appID", branch.ApplicationID).
		WithField("branchID", branch.ID).
		Info("branch created")

	return &domain.Branch{
		ID:            branch.ID,
		ApplicationID: branch.ApplicationID,
		BranchName:    branch.BranchName,
		BuildType:     builder.BuildTypeFromString(branch.BuildType),
	}, nil
}

func (r *applicationRepository) GetBranchByID(ctx context.Context, id string) (*domain.Branch, error) {
	const errMsg = "failed to GetBranchByID: %w"

	branch, err := models.FindBranch(ctx, r.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	return &domain.Branch{
		ID:            branch.ID,
		ApplicationID: branch.ApplicationID,
		BranchName:    branch.BranchName,
		BuildType:     builder.BuildTypeFromString(branch.BuildType),
	}, nil
}

func (r *applicationRepository) GetBranchByRepoAndBranchName(ctx context.Context, repoURL string, branch string) (*domain.Branch, error) {
	const errMsg = "failed to GetBranchByRepoAndBranchName: %w"

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
		qm.Load(models.ApplicationRels.Branches, models.BranchWhere.BranchName.EQ(branch)),
		models.ApplicationWhere.DeletedAt.IsNull(),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	if len(app.R.Branches) > 0 {
		branch := app.R.Branches[0]
		return &domain.Branch{
			ID:            branch.ID,
			ApplicationID: branch.ApplicationID,
			BranchName:    branch.BranchName,
			BuildType:     builder.BuildTypeFromString(branch.BuildType),
		}, nil
	}
	return nil, ErrNotFound
}

func (r *applicationRepository) SetWebsite(ctx context.Context, branchID string, fqdn string, httpPort int) error {
	const errMsg = "failed to SetWebsite: %w"

	branch, err := models.Branches(
		qm.Load(models.BranchRels.Website),
		models.BranchWhere.ID.EQ(branchID),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf(errMsg, err)
	}

	ws := branch.R.Website
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
	if err := branch.SetWebsite(ctx, r.db, true, ws); err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}
