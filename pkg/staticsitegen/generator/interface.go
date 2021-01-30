package generator

import "github.com/traPtitech/neoshowcase/pkg/storage"

type Engine interface {
	Init(s storage.Storage) error
	Reconcile(sites []*Site) error
}

type Site struct {
	ID            string
	FQDN          string
	ArtifactID    string
	ApplicationID string
}
