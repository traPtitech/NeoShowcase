package oidc

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func NewGitHubOIDCProvider(clientID, clientSecret string) (*oauth2.Config, *oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(context.TODO(), "https://github.com/")
	if err != nil {
		return nil, nil, err
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		Scopes: []string{string(github.ScopePublicRepo), string(github.ScopeUser), string(github.ScopeDeleteRepo)},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	return config, verifier, nil
}
