package dockerimpl

import (
	"bytes"
	"context"

	docker "github.com/fsouza/go-dockerclient"
)

type LogOptions struct {
	Tail   string
	Stdout bool
	Stderr bool
	Since  int64
}

func (b *dockerBackend) GetContainerLog(ctx context.Context, appID string, envID string, opt LogOptions) (string, error) {
	str := &bytes.Buffer{}
	logopts := docker.LogsOptions{
		Context:      ctx,
		Container:    containerName(appID, envID),
		OutputStream: str,
		Tail:         opt.Tail,
		Stdout:       opt.Stdout,
		Stderr:       opt.Stderr,
		Since:        opt.Since,
		RawTerminal:  true,
	}
	err := b.c.Logs(logopts)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}
