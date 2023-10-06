package builder

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/regclient/regclient/types/ref"
)

func (s *builderService) buildStaticExtract(
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
	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"local-src": st.repositoryTempDir,
		},
		Session: s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildStaticCleanup(
	ctx context.Context,
	st *state,
) error {
	tagRef, err := ref.New(s.tmpDestImage(st.app, st.build))
	if err != nil {
		return err
	}
	return s.registry.TagDelete(ctx, tagRef)
}
