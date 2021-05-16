package usecase

import (
	"context"
	"database/sql"
	"fmt"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type VirtualHostService interface {
	GenerateRandomDomain(ctx context.Context) (string, error)
}

type virtualHostService struct {
	db *sql.DB
}

func NewVirtualHostService(db *sql.DB) VirtualHostService {
	return &virtualHostService{db: db}
}

func (s *virtualHostService) GenerateRandomDomain(ctx context.Context) (string, error) {
	d, err := models.AvailableDomains(
		models.AvailableDomainWhere.Subdomain.EQ(true),
	).One(context.Background(), s.db)
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
		ng, err := models.Websites(models.WebsiteWhere.FQDN.EQ(candidate)).Exists(context.Background(), s.db)
		if err != nil {
			return "", fmt.Errorf("db error: %w", err)
		}
		if !ng {
			break
		}
	}

	return candidate, nil
}
