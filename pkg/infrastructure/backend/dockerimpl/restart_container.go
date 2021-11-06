package dockerimpl

import (
	"context"
	"fmt"
)

func (b *dockerBackend) RestartContainer(ctx context.Context, appID string, branchID string) error {
	// コンテナを止めて5秒待ったのちkillし, コンテナを再起動
	err := b.c.RestartContainer(containerName(appID, branchID), timeout)
	if err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}
	return nil
}
