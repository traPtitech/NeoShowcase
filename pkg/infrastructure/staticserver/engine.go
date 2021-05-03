package staticserver

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type WebServerDocumentRootPath string

type WebServerPort int

type Engine interface {
	Start(ctx context.Context) error
	Reconcile(sites []*domain.Site) error
	Shutdown(ctx context.Context) error
}
