package dockerimpl

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) unregisterIngress(_ context.Context, _ *domain.Application, website *domain.Website) error {
	err := os.Remove(b.configFile(website.FQDN))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
