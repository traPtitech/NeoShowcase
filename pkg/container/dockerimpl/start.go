package dockerimpl

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/container"
)

func (m *Manager) Start(ctx context.Context, args container.StartArgs) (*container.StartResult, error) {
	err := m.c.StartContainerWithContext(containerName(args.ApplicationID, args.EnvironmentID), nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}
	return &container.StartResult{}, nil
}
