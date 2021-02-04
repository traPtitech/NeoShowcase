package appmanager

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/idgen"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type envImpl struct {
	m       *managerImpl
	dbmodel *models.Environment
}

func (env *envImpl) GetID() string {
	return env.dbmodel.ID
}

func (env *envImpl) GetBranchName() string {
	return env.dbmodel.BranchName
}

func (env *envImpl) GetBuildType() BuildType {
	return BuildTypeFromString(env.dbmodel.BuildType)
}

func (env *envImpl) SetupWebsite(fqdn string, httpPort int) error {
	if httpPort <= 0 || httpPort > 65535 {
		httpPort = 80
	}
	if fqdn == "" {
		// ドメインを自動的に生成
		d, err := env.m.generateRandomDomain()
		if err != nil {
			return fmt.Errorf("failed to generateRandomDomain: %w", err)
		}
		fqdn = d
	}

	ws := env.dbmodel.R.Website
	if ws != nil {
		// テーブルの情報を更新
		ws.FQDN = fqdn
		ws.HTTPPort = httpPort
		if _, err := ws.Update(context.Background(), env.m.db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to update website: %w", err)
		}
		return nil
	}

	// Websiteをテーブルに挿入
	ws = &models.Website{
		ID:       idgen.New(),
		FQDN:     fqdn,
		HTTPPort: httpPort,
	}
	if err := env.dbmodel.SetWebsite(context.Background(), env.m.db, true, ws); err != nil {
		return fmt.Errorf("failed to SetWebsite: %w", err)
	}
	return nil
}
