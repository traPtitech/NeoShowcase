package oidc

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func NewGitHubOIDCProvider(clientID, clientSecret string) (*oauth2.Config, *oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(context.TODO(), "")
	if err != nil {
		return nil, nil, err
	}

	config := &oauth2.Config{
		// TODO: fill
	}

	verifier := provider.Verifier(&oidc.Config{
		// TODO: fill
	})

	return config, verifier, nil
}
