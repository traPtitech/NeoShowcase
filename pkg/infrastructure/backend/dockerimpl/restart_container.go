package dockerimpl

import (
	"context"
	"fmt"
)

func (b *dockerBackend) RestartContainer(ctx context.Context, appID string, envID string) error {
	// コンテナを止めて5秒待ったのちkillし, コンテナを再起動
	err := b.c.RestartContainer(containerName(appID, envID), timeout)
	if err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}
	return nil
}
