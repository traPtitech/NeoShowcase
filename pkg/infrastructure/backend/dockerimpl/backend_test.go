package dockerimpl

import (
	"context"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

func prepareManager(t *testing.T) (*Backend, *client.Client) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DOCKER_TESTS")); !ok {
		t.SkipNow()
	}

	// Dockerデーモンに接続 (DooD)
	c, err := NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	config := Config{}
	config.ConfDir = "../../../../.local-dev/traefik"
	config.Network = "neoshowcase_apps"
	config.Routing.Type = routingTypeTraefik
	m, err := NewDockerBackend(c, config, builder.ImageConfig{})
	require.NoError(t, err)
	err = m.Start(context.Background())
	require.NoError(t, err)

	res, err := c.ImagePull(context.Background(), "alpine:latest", types.ImagePullOptions{})
	require.NoError(t, err)
	_, err = io.ReadAll(res)
	require.NoError(t, err)
	require.NoError(t, res.Close())

	return m, c
}
