package builder

import (
	"context"
	"io"
)

type BuildpackBackend interface {
	Pack(ctx context.Context, repoDir string, logWriter io.Writer, imageDest string) (path string, err error)
}
