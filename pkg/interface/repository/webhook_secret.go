package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

type WebhookSecretRepository interface {
	// GetWebhookSecretKeys 指定したリポジトリのWebhookシークレットキーを取得します
	GetWebhookSecretKeys(ctx context.Context, repositoryUrl string) ([]string, error)
}

type webhookSecretRepository struct {
	db *sql.DB
}

func NewWebhookSecretRepository(db *sql.DB) WebhookSecretRepository {
	return &webhookSecretRepository{db: db}
}

func (r *webhookSecretRepository) GetWebhookSecretKeys(ctx context.Context, repositoryUrl string) ([]string, error) {
	_, err := models.Repositories(models.RepositoryWhere.Remote.EQ(repositoryUrl)).
		One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetWebhookSecretKeys: %w", err)
	}

	// TODO 実装
	return []string{}, nil
}
