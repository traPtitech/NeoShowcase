package webserver

import (
	"context"

	storage2 "github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

type Engine interface {
	Init(s storage2.Storage) error
	Start(ctx context.Context) error
	Reconcile(sites []*Site) error
	Close(ctx context.Context) error
}

type Site struct {
	ID            string
	FQDN          string
	ArtifactID    string
	ApplicationID string
}
