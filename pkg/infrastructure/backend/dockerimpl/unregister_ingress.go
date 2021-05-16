package dockerimpl

import (
	"context"
	"os"
	"path/filepath"

	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *dockerBackend) UnregisterIngress(ctx context.Context, appID string, envID string) error {
	conf := filepath.Join(b.ingressConfDir, containerName(appID, envID)+".yaml")
	if util.FileExists(conf) {
		return os.Remove(conf)
	}
	return nil
}
