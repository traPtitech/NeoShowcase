package dockerimpl

import (
	"context"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

const key = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMwwAAAJCX+efxl/nn
8QAAAAtzc2gtZWQyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMww
AAAEA+FzwWKIYduEDOqkEOZ2wmxZWPc2wpZeWj+J8e3Q6x0QLWIALnhPWigJA3njldxcM9
1PEkQQN231Ow7czCg8zDAAAADG1vdG9AbW90by13cwE=
-----END OPENSSH PRIVATE KEY-----`

func prepareManager(t *testing.T) (*dockerBackend, *client.Client) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DOCKER_TESTS")); !ok {
		t.SkipNow()
	}

	// Dockerデーモンに接続 (DooD)
	c, err := NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	sshKey, err := ssh2.NewPublicKeys("", []byte(key), "")
	require.NoError(t, err)
	m, err := NewDockerBackend(c, Config{
		SSH:     sshConfig{Port: 2201},
		ConfDir: "../../../../.local-dev/traefik",
		Network: "neoshowcase_apps",
	}, sshKey, builder.ImageConfig{}, nil, nil)
	require.NoError(t, err)
	err = m.Start(context.Background())
	require.NoError(t, err)

	res, err := c.ImagePull(context.Background(), "alpine:latest", types.ImagePullOptions{})
	require.NoError(t, err)
	_, err = io.ReadAll(res)
	require.NoError(t, err)
	require.NoError(t, res.Close())

	return m.(*dockerBackend), c
}
