package builder

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/tonistiigi/fsutil"
)

func (s *ServiceImpl) buildStaticExtract(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	ls := llb.Image(s.tmpDestImage(st.app, st.build))
	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(ls, st.staticDest, "/", &llb.CopyInfo{
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

func (s *ServiceImpl) buildStaticCleanup(
	ctx context.Context,
	st *state,
) error {
	return s.regclient.DeleteImage(ctx, s.imageConfig.TmpImageName(st.app.ID), s.imageTag(st.build))
}
