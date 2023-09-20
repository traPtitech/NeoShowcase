package builder

import (
	"context"
	"io"
	"path/filepath"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/regclient/regclient/types/ref"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (s *builderService) buildRuntimeBuildpack(
	ctx context.Context,
	st *state,
	bc *domain.BuildConfigRuntimeBuildpack,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	buildDir := filepath.Join(st.repositoryTempDir, contextDir)
	env, err := s.appEnv(ctx, st.app)
	if err != nil {
		return err
	}
	_, err = s.buildpack.Pack(ctx, buildDir, s.destImage(st.app, st.build), env, st.Logger())
	return err
}

func (s *builderService) buildStaticBuildpackPack(
	ctx context.Context,
	st *state,
	bc *domain.BuildConfigStaticBuildpack,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	buildDir := filepath.Join(st.repositoryTempDir, contextDir)
	env, err := s.appEnv(ctx, st.app)
	if err != nil {
		return err
	}
	path, err := s.buildpack.Pack(ctx, buildDir, s.tmpDestImage(st.app, st.build), env, st.Logger())
	if err != nil {
		return err
	}
	st.buildpackDest = path
	return nil
}

func (s *builderService) buildStaticBuildpackExtract(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticBuildpack,
) error {
	ls := llb.Image(s.tmpDestImage(st.app, st.build))
	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(ls, filepath.Join(st.buildpackDest, bc.ArtifactPath), "/", &llb.CopyInfo{
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

func (s *builderService) buildStaticBuildpackCleanup(
	ctx context.Context,
	st *state,
) error {
	tagRef, err := ref.New(s.config.TmpImageName(st.app.ID) + ":" + st.build.ID)
	if err != nil {
		return err
	}
	return s.registry.TagDelete(ctx, tagRef)
}
