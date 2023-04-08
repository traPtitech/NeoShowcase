package domain

import (
	"context"
)

type (
	StaticServerDocumentRootPath string
	StaticServerPort             int
)

type SSEngine interface {
	Start(ctx context.Context) error
	Reconcile(sites []*StaticSite) error
	Shutdown(ctx context.Context) error
}

type StaticSite struct {
	Application *Application
	Website     *Website
	ArtifactID  string
}
