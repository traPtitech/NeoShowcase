package dockerimpl

import (
	"bytes"
	"context"

	docker "github.com/fsouza/go-dockerclient"
)

type LogsOptions struct {
	Tail  string
	Since int64
}

func (b *dockerBackend) GetContainerStdOut(ctx context.Context, appID string, envID string, opt LogsOptions) (string, error) {
	str := &bytes.Buffer{}
	logopts := docker.LogsOptions{
		Context:      ctx,
		Container:    containerName(appID, envID),
		OutputStream: str,
		Tail:         opt.Tail,
		Stdout:       true,
		Since:        opt.Since,
	}
	err := b.c.Logs(logopts)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}
