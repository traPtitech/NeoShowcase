package webserver

import (
	"context"
	"github.com/traPtitech/neoshowcase/pkg/storage"
)

type Engine interface {
	Init(s storage.Storage) error
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
