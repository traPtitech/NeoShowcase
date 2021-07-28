package domain

import (
	"context"
)

type WebServerDocumentRootPath string

type WebServerPort int

type Engine interface {
	Start(ctx context.Context) error
	Reconcile(sites []*Site) error
	Shutdown(ctx context.Context) error
}

type Site struct {
	ID            string
	FQDN          string
	ArtifactID    string
	ApplicationID string
}
