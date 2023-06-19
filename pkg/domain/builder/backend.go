package builder

import (
	"context"
	"io"
)

type BuildpackBackend interface {
	Pack(
		ctx context.Context,
		repoDir string,
		imageDest string,
		env map[string]string,
		logWriter io.Writer,
	) (path string, err error)
}
