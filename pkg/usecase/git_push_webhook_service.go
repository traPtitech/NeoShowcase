package usecase

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type GitPushWebhookService interface {
	VerifySignature(ctx context.Context, repoURL string, signature string, body []byte) (valid bool, err error)
	CheckRepositoryExists(ctx context.Context, args CheckRepositoryExistsArgs) (exists bool, err error)
}

type gitPushWebhookService struct {
	gitRepo    repository.GitrepositoryRepository
	secretRepo repository.WebhookSecretRepository
}

type CheckRepositoryExistsArgs struct {
	ProviderID string
	Owner      string
	RepoName   string
}

func NewGitPushWebhookService(secretRepo repository.WebhookSecretRepository, gitRepo repository.GitrepositoryRepository) GitPushWebhookService {
	return &gitPushWebhookService{
		secretRepo: secretRepo,
		gitRepo:    gitRepo,
	}
}

func (s *gitPushWebhookService) VerifySignature(ctx context.Context, repoURL string, signature string, body []byte) (bool, error) {
	secrets, err := s.secretRepo.GetWebhookSecretKeys(ctx, repoURL)
	if err != nil {
		return false, fmt.Errorf("failed to GetWebhookSecretKeys: %w", err)
	}

	// シグネチャの検証(secretsのうちひとつでもtrueならOK)
	valid := false
	for _, secret := range secrets {
		valid = valid || domain.VerifySignature(sha256.New, body, []byte(secret), signature)
	}
	return valid, nil
}

func (s *gitPushWebhookService) CheckRepositoryExists(ctx context.Context, args CheckRepositoryExistsArgs) (bool, error) {
	r := repository.GetRepositoryArgs{
		ProviderID: args.ProviderID,
		Owner:      args.Owner,
		Name:       args.RepoName,
	}
	_, err := s.gitRepo.GetRepository(ctx, r)

	if err != nil {
		if err == repository.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to GetRepository: %w", err)
	}

	return true, nil

}
