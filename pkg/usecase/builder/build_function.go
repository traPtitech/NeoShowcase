package builder

import (
	"context"
	"path"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/tonistiigi/fsutil"
)

func (s *ServiceImpl) extractFunctionArtifact(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	artifactPath := st.artifactTempFile.Name()
	artifactName := path.Base(artifactPath)
	artifactDir := path.Dir(artifactPath)

	ls := llb.Image(s.tmpDestImage(st.app, st.build))
	// ビルドで生成された単一ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(ls, st.functionDest, "/"+artifactName, &llb.CopyInfo{
			CreateDestPath: true,
			AllowWildcard:  true,
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
			Type:      buildkit.ExporterLocal,
			OutputDir: artifactDir,
		}},
		LocalMounts: map[string]fsutil.FS{
			"local-src": mount,
		},
		Session: s.authSessions(),
	}, ch)
	return err
}
