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
	// ReclaimOrphans deletes resources left behind by builds of a previous process of this builder,
	// such as when the builder crashed in the middle of a build.
	// It must be called before the first Pack, since it also deletes resources of builds running at the time.
	ReclaimOrphans(ctx context.Context) error
	// ReclaimEnded deletes resources left behind by builds that have ended, whichever builder ran them.
	// Unlike ReclaimOrphans, it leaves running builds alone and may be called at any time.
	ReclaimEnded(ctx context.Context) error
}
