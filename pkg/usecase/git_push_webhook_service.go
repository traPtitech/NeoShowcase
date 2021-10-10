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
}

type gitPushWebhookService struct {
	repository.WebhookSecretRepository
}

func NewGitPushWebhookService(repo repository.WebhookSecretRepository) GitPushWebhookService {
	return &gitPushWebhookService{WebhookSecretRepository: repo}
}

func (s *gitPushWebhookService) VerifySignature(ctx context.Context, repoURL string, signature string, body []byte) (bool, error) {
	secrets, err := s.GetWebhookSecretKeys(ctx, repoURL)
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
