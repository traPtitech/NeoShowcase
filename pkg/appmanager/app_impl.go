package appmanager

import (
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/idgen"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *managerImpl) GetApp(appID string) (App, error) {
	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(qm.Rels(models.ApplicationRels.Environments, models.EnvironmentRels.Website)),
		models.ApplicationWhere.DeletedAt.IsNull(),
		models.ApplicationWhere.ID.EQ(appID),
	).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetApp: %w", err)
	}

	return &appImpl{
		m:       m,
		dbmodel: app,
	}, nil
}

func (m *managerImpl) GetAppByRepository(repo string) (App, error) {
	repoModel, err := models.Repositories(models.RepositoryWhere.Remote.EQ(repo)).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetAppByRepository: %w", err)
	}

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(qm.Rels(models.ApplicationRels.Environments, models.EnvironmentRels.Website)),
		models.ApplicationWhere.DeletedAt.IsNull(),
		models.ApplicationWhere.RepositoryID.EQ(repoModel.ID),
	).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetApp: %w", err)
	}

	return &appImpl{
		m:       m,
		dbmodel: app,
	}, nil
}

func (m *managerImpl) GetAppByEnvironment(envID string) (App, error) {
	env, err := models.FindEnvironment(context.Background(), m.db, envID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetAppByEnvironment: %w", err)
	}

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(qm.Rels(models.ApplicationRels.Environments, models.EnvironmentRels.Website)),
		models.ApplicationWhere.DeletedAt.IsNull(),
		models.ApplicationWhere.RepositoryID.EQ(env.ID),
	).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetApp: %w", err)
	}

	return &appImpl{
		m:       m,
		dbmodel: app,
	}, nil
}

func (m *managerImpl) CreateApp(args CreateAppArgs) (App, error) {
	// リポジトリ情報を設定
	repo, err := models.Repositories(models.RepositoryWhere.Remote.EQ(args.RepositoryURL)).One(context.Background(), m.db)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	} else if repo == nil {
		repo = &models.Repository{
			ID:     idgen.New(),
			Remote: args.RepositoryURL,
		}
		if err := repo.Insert(context.Background(), m.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf("failed to insert repository: %w", err)
		}
	}

	// アプリケーション作成
	app := &models.Application{
		ID:           idgen.New(),
		Owner:        args.Owner,
		Name:         args.Name,
		RepositoryID: repo.ID,
	}
	if err := app.Insert(context.Background(), m.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to insert application: %w", err)
	}
	if err := app.SetRepository(context.Background(), m.db, false, repo); err != nil {
		return nil, fmt.Errorf("failed to associate repository: %w", err)
	}

	log.WithField("appID", app.ID).Info("app created")
	appI := &appImpl{
		m:       m,
		dbmodel: app,
	}

	// 初期Env作成
	if _, err := appI.CreateEnv(args.BranchName, args.BuildType); err != nil {
		return nil, err
	}
	return appI, nil
}

type appImpl struct {
	m *managerImpl
	// dbmodelはEnvironments.WebsiteとRepositoryがプレロードされてる事
	dbmodel *models.Application
}

func (app *appImpl) GetID() string {
	return app.dbmodel.ID
}

func (app *appImpl) GetName() string {
	return app.dbmodel.Name
}

func (app *appImpl) GetEnvs() []Env {
	result := make([]Env, 0)
	for _, env := range app.dbmodel.R.Environments {
		result = append(result, &envImpl{
			m:       app.m,
			dbmodel: env,
		})
	}
	return result
}

func (app *appImpl) CreateEnv(branchName string, buildType BuildType) (Env, error) {
	// 指定したブランチの環境が存在しないことを確認
	for _, env := range app.dbmodel.R.Environments {
		if env.BranchName == branchName {
			return nil, fmt.Errorf("the environment for branch `%s` has already existed", branchName)
		}
	}

	env := &models.Environment{
		ID:         idgen.New(),
		BranchName: branchName,
		BuildType:  buildType.String(),
	}
	if err := app.dbmodel.AddEnvironments(context.Background(), app.m.db, true, env); err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}
	log.WithField("appID", env.ApplicationID).WithField("envID", env.ID).Info("env created")
	return &envImpl{m: app.m, dbmodel: env}, nil
}

func (app *appImpl) Start(args AppStartArgs) error {
	var env *models.Environment
	for _, _env := range app.dbmodel.R.Environments {
		if _env.ID == args.EnvironmentID {
			env = _env
			break
		}
	}
	if env == nil {
		return fmt.Errorf("environtment not found: %s", args.EnvironmentID)
	}

	switch env.BuildType {
	case models.EnvironmentsBuildTypeImage:
		if args.BuildID == "" {
			// buildIDの指定がない場合は最新のビルドを使用
			build, err := env.BuildLogs(
				qm.OrderBy(fmt.Sprintf("%s DESC", models.BuildLogColumns.StartedAt)),
				models.BuildLogWhere.Result.EQ(models.BuildLogsResultSUCCEEDED),
			).One(context.Background(), app.m.db)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("no successful build exists")
				}
				return fmt.Errorf("failed to get BuildLogs: %w", err)
			}
			args.BuildID = build.ID
		} else {
			// buildIDのビルドが存在するかどうか確認
			ok, err := models.BuildLogExists(context.Background(), app.m.db, args.BuildID)
			if err != nil {
				return fmt.Errorf("failed to BuildLogExists: %w", err)
			}
			if !ok {
				return fmt.Errorf("build (%s) was not found", args.BuildID)
			}
		}

		// HTTP公開設定があれば取得
		var httpProxy *container.HTTPProxy
		website, err := env.Website().One(context.Background(), app.m.db)
		if err == nil {
			httpProxy = &container.HTTPProxy{
				Domain: website.FQDN,
				Port:   website.HTTPPort,
			}
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("failed to query website: %w", err)
		}

		_, err = app.m.cm.Create(context.Background(), container.CreateArgs{
			ApplicationID: app.GetID(),
			EnvironmentID: env.ID,
			ImageName:     app.m.getFullImageName(app),
			ImageTag:      args.BuildID,
			HTTPProxy:     httpProxy,
			Recreate:      true,
		})
		if err != nil {
			return fmt.Errorf("failed to Create container: %w", err)
		}

		env.BuildID = null.StringFrom(args.BuildID)
		if _, err := env.Update(context.Background(), app.m.db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to Update website: %w", err)
		}

	case models.EnvironmentsBuildTypeStatic:
		// TODO 実装
		log.Fatalf("NOT IMPLEMENTED")

	default:
		return fmt.Errorf("unknown build type: %s", env.BuildType)
	}
	return nil
}

func (app *appImpl) RequestBuild(ctx context.Context, envID string) error {
	var env *models.Environment
	for _, _env := range app.dbmodel.R.Environments {
		if _env.ID == envID {
			env = _env
			break
		}
	}
	if env == nil {
		return fmt.Errorf("environtment not found: %s", envID)
	}

	switch env.BuildType {
	case models.EnvironmentsBuildTypeImage:
		_, err := app.m.builder.StartBuildImage(ctx, &builderApi.StartBuildImageRequest{
			ImageName: app.m.getImageName(app),
			Source: &builderApi.BuildSource{
				RepositoryUrl: app.dbmodel.R.Repository.Remote, // TODO ブランチ・タグ指定に対応
			},
			Options:       &builderApi.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			EnvironmentId: env.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build image: %w", err)
		}

	case models.EnvironmentsBuildTypeStatic:
		_, err := app.m.builder.StartBuildStatic(ctx, &builderApi.StartBuildStaticRequest{
			Source: &builderApi.BuildSource{
				RepositoryUrl: app.dbmodel.R.Repository.Remote, // TODO ブランチ・タグ指定に対応
			},
			Options:       &builderApi.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			EnvironmentId: env.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build static: %w", err)
		}

	default:
		return fmt.Errorf("unknown build type: %s", env.BuildType)
	}
	return nil
}
