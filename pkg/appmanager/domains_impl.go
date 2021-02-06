package appmanager

import (
	"context"
	"database/sql"
	"fmt"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/traPtitech/neoshowcase/pkg/models"
)

// generateRandomDomain ランダムなドメインを生成する
func (m *managerImpl) generateRandomDomain() (string, error) {
	d, err := models.AvailableDomains(
		models.AvailableDomainWhere.Subdomain.EQ(true),
	).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			// DBに親ドメインが登録されてないと生成できない
			return "", fmt.Errorf("auto domain generation is not available")
		}
		return "", err
	}

	// 被らないものを生成
	var candidate string
	for {
		candidate = petname.Generate(3, "-") + "." + d.Domain
		ng, err := models.Websites(models.WebsiteWhere.FQDN.EQ(candidate)).Exists(context.Background(), m.db)
		if err != nil {
			return "", fmt.Errorf("db error: %w", err)
		}
		if !ng {
			break
		}
	}

	return candidate, nil
}
