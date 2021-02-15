package appmanager

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/models"
)

func (m *managerImpl) GetRepoByURL(url string) (Repo, error) {
	repoModel, err := models.Repositories(models.RepositoryWhere.Remote.EQ(url)).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetRepoByURL: %w", err)
	}
	return &repoImpl{
		m:       m,
		dbModel: repoModel,
	}, nil
}

type repoImpl struct {
	m       *managerImpl
	dbModel *models.Repository
}

func (r *repoImpl) GetID() string {
	return r.dbModel.ID
}

func (r *repoImpl) GetGitURL() string {
	return r.dbModel.Remote
}

func (r *repoImpl) GetWebhookSecret() string {
	panic("implement me") // TODO
}

func (r *repoImpl) SetWebhookSecret(secret string) error {
	panic("implement me") // TODO
}
