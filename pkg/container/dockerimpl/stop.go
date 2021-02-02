package dockerimpl

import (
	"context"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/container"
)

func (m *Manager) Stop(ctx context.Context, args container.StopArgs) (*container.StopResult, error) {
	// コンテナを止めて5秒待ったのちkill
	err := m.c.StopContainer(containerName(args.ApplicationID, args.EnvironmentID), timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to stop container: %w", err)
	}
	return &container.StopResult{}, nil
}
