package dockerimpl

import (
	"os"
	"strconv"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

func prepareManager(t *testing.T, bus eventbus.Bus) (*dockerBackend, *docker.Client) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DOCKER_TESTS")); !ok {
		t.SkipNow()
	}

	// Dockerデーモンに接続 (DooD)
	c, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	m, err := NewDockerBackend(c, bus, "./local-dev/traefik")
	if err != nil {
		t.Fatal(err)
	}

	require.NoError(t, c.PullImage(docker.PullImageOptions{
		Repository: "alpine",
		Tag:        "latest",
	}, docker.AuthConfiguration{}))

	return m.(*dockerBackend), c
}
