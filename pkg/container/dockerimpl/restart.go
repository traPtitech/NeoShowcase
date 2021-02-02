package dockerimpl

import (
	"context"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/container"
)

func (m *Manager) Restart(ctx context.Context, args container.RestartArgs) (*container.RestartResult, error) {
	// コンテナを止めて5秒待ったのちkillし, コンテナを再起動
	err := m.c.RestartContainer(containerName(args.ApplicationID, args.EnvironmentID), timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to restart container: %w", err)
	}
	return &container.RestartResult{}, nil
}
