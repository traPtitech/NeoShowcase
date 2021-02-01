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
	switch app.dbmodel.BuildType {
	case models.ApplicationsBuildTypeImage:
		if args.BuildID == "" {
			// buildIDの指定がない場合は最新のビルドを使用
			build, err := app.dbmodel.BuildLogs(
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
		website, err := app.dbmodel.Website().One(context.Background(), app.m.db)
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

		if website != nil {
			website.BuildID = null.StringFrom(args.BuildID)
			if _, err := website.Update(context.Background(), app.m.db, boil.Infer()); err != nil {
				return fmt.Errorf("failed to Update website: %w", err)
			}
		}

	case models.ApplicationsBuildTypeStatic:
		// TODO 実装
		log.Fatalf("NOT IMPLEMENTED")

	default:
		log.Fatalf("unknown build type: %s", app.dbmodel.BuildType)
	}
	return nil
}

// RequestBuild builderにappのビルドをリクエストする
func (app *appImpl) RequestBuild(ctx context.Context) error {
	switch app.dbmodel.BuildType {
	case models.ApplicationsBuildTypeImage:
		_, err := app.m.builder.StartBuildImage(ctx, &builderApi.StartBuildImageRequest{
			ImageName: app.m.getImageName(app),
			Source: &builderApi.BuildSource{
				RepositoryUrl: app.dbmodel.R.Repository.Remote, // TODO ブランチ・タグ指定に対応
			},
			Options:       &builderApi.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			ApplicationId: app.GetID(),
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build image: %w", err)
		}
		return nil

	case models.ApplicationsBuildTypeStatic:
		_, err := app.m.builder.StartBuildStatic(ctx, &builderApi.StartBuildStaticRequest{
			Source: &builderApi.BuildSource{
				RepositoryUrl: app.dbmodel.R.Repository.Remote, // TODO ブランチ・タグ指定に対応
			},
			Options:       &builderApi.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			ApplicationId: app.GetID(),
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build static: %w", err)
		}
		return nil

	default:
		log.Fatalf("unknown build type: %s", app.dbmodel.BuildType)
	}
	return nil
}
