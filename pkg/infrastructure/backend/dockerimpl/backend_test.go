package dockerimpl

import (
	"context"
	"os"
	"strconv"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

func prepareManager(t *testing.T) (*dockerBackend, *docker.Client) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DOCKER_TESTS")); !ok {
		t.SkipNow()
	}

	// Dockerデーモンに接続 (DooD)
	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	m, err := NewDockerBackend(c, Config{
		ConfDir: "../../../../.local-dev/traefik",
		Network: "neoshowcase_apps",
	}, builder.ImageConfig{})
	require.NoError(t, err)
	err = m.Start(context.Background())
	require.NoError(t, err)

	require.NoError(t, c.PullImage(docker.PullImageOptions{
		Repository: "alpine",
		Tag:        "latest",
	}, docker.AuthConfiguration{}))

	return m.(*dockerBackend), c
}
