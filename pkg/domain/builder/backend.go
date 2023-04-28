package builder

import (
	"context"
	"io"
)

type BuildpackBackend interface {
	Pack(ctx context.Context, repoDir string, imageDest string, logWriter io.Writer) error
}
