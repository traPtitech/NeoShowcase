package builder

import (
	"context"
	"io"
)

type BuildpackBackend interface {
	// Pack builds the app in repoDir into imageDest, and returns the path of the app directory inside the image.
	Pack(
		ctx context.Context,
		buildID string,
		repoDir string,
		imageDest string,
		imageConfig ImageConfig,
		env map[string]string,
		logWriter io.Writer,
	) (path string, err error)
}
