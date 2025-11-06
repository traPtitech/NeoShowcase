package builder

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/tonistiigi/fsutil"
)

func (s *ServiceImpl) buildExtractFolderToTar(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	path string,
) error {
	ls := llb.Image(s.tmpDestImage(st.app, st.build))
	def, err := llb.
		Scratch().
		File(llb.Copy(ls, path, "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to marshal llb")
	}
	mount, err := fsutil.NewFS(st.repositoryTempDir)
	if err != nil {
		return errors.Wrap(err, "invalid mount dir")
	}
	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalMounts: map[string]fsutil.FS{
			"local-src": mount,
		},
		Session: s.authSessions(),
	}, ch)
	return err
}

func (s *ServiceImpl) buildStaticExtract(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	return s.buildExtractFolderToTar(ctx, st, ch, st.staticDest)
}

func (s *ServiceImpl) extractFunctionArtifact(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	return s.buildExtractFolderToTar(ctx, st, ch, st.functionDest)
}
