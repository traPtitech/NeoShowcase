package usecase

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type GitPushWebhookService interface {
	VerifySignature(ctx context.Context, repoURL string, signature string, body []byte) (valid bool, err error)
	CheckRepositoryExists(ctx context.Context, repoURL, owner, name string) (bool, error)
}

type gitPushWebhookService struct {
	repo repository.GitRepositoryRepository
}

var ErrProviderNotFound = fmt.Errorf("provider not found")

func NewGitPushWebhookService(repo repository.GitRepositoryRepository) GitPushWebhookService {
	return &gitPushWebhookService{
		repo: repo,
	}
}

func (s *gitPushWebhookService) VerifySignature(ctx context.Context, repoURL string, signature string, body []byte) (bool, error) {
	u, err := url.Parse(repoURL)

	if err != nil {
		return false, err
	}

	prov, err := s.repo.GetProviderByHost(ctx, u.Host)
	if err != nil {
		return false, ErrProviderNotFound
	}

	// シグネチャの検証
	valid := domain.VerifySignature(sha256.New, body, []byte(prov.Secret), signature)

	return valid, nil
}

func (s *gitPushWebhookService) CheckRepositoryExists(ctx context.Context, repoURL, owner, name string) (bool, error) {
	u, err := url.Parse(repoURL)

	if err != nil {
		return false, err
	}

	prov, err := s.repo.GetProviderByHost(ctx, u.Host)

	if err != nil {
		return false, ErrProviderNotFound
	}
	r := repository.GetRepositoryArgs{
		ProviderID: prov.ID,
		Owner:      owner,
		Name:       name,
	}

	_, err = s.repo.GetRepository(ctx, r)

	if err != nil {
		if err == repository.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to GetRepository: %w", err)
	}

	return true, nil

}
