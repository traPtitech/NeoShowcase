package appmanager

import (
	"context"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type appImpl struct {
	m       *managerImpl
	dbmodel *models.Application
}

func (app *appImpl) GetID() string {
	return app.dbmodel.ID
}

func (app *appImpl) GetName() string {
	return app.dbmodel.Name
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

		// TODO 既に存在する場合はCreateせずにStartするようにする
		_, err = app.m.cm.Create(context.Background(), container.CreateArgs{
			ApplicationID: app.GetID(),
			ImageName:     app.m.getFullImageName(app),
			ImageTag:      args.BuildID,
			HTTPProxy:     httpProxy,
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

// RequestBuild builderにappのビルドをリクエストする
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
