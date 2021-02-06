package dockerimpl

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func prepareManager(t *testing.T) (*Manager, *docker.Client, *hub.Hub) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DOCKER_TESTS")); !ok {
		t.SkipNow()
	}
	bus := hub.New()

	m, err := NewManager(bus)
	if err != nil {
		t.Fatal(err)
	}

	require.NoError(t, m.c.PullImage(docker.PullImageOptions{
		Repository: "alpine",
		Tag:        "latest",
	}, docker.AuthConfiguration{}))

	return m, m.c, bus
}

func TestNewManager(t *testing.T) {
	_, c, _ := prepareManager(t)

	_, err := c.Version()
	assert.NoError(t, err)
}
