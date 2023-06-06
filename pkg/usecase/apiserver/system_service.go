package apiserver

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/motoki317/sc"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (s *Service) GetSystemPublicKey(_ context.Context) string {
	encoded := domain.Base64EncodedPublicKey(s.pubKey.Signer.PublicKey())
	return encoded + " neoshowcase"
}

func (s *Service) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	return s.controller.GetAvailableDomains(ctx)
}

func (s *Service) GetAvailablePorts(ctx context.Context) (domain.AvailablePortSlice, error) {
	available, err := s.controller.GetAvailablePorts(ctx)
	if err != nil {
		return nil, err
	}
	return available, nil
}

type tmpKeyPairService struct {
	*sc.Cache[string, ed25519.PrivateKey]
}

func newTmpKeyPairService() *tmpKeyPairService {
	return &tmpKeyPairService{
		Cache: sc.NewMust(func(ctx context.Context, key string) (ed25519.PrivateKey, error) {
			_, privKey, err := ed25519.GenerateKey(nil)
			if err != nil {
				return nil, err
			}
			return privKey, nil
		}, 1*time.Hour, 1*time.Hour, sc.WithCleanupInterval(1*time.Hour)),
	}
}

func (s *Service) GenerateKeyPair(ctx context.Context) (keyID string, publicKey string, err error) {
	keyID = domain.NewID()
	privKey, err := s.tmpKeys.Get(ctx, keyID)
	if err != nil {
		return "", "", errors.Wrap(err, "generating ed25519 key")
	}
	pubKey, err := ssh.NewPublicKey(privKey.Public())
	if err != nil {
		return "", "", errors.Wrap(err, "creating public key")
	}
	encoded := domain.Base64EncodedPublicKey(pubKey)
	return keyID, encoded + " neoshowcase", nil
}
