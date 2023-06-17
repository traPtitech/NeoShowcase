package builder

import (
	"context"
	"io"
	"path/filepath"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
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
	_, err := s.buildpack.Pack(ctx, buildDir, st.Logger(), s.destImage(st.app, st.build))
	return err
}

func (s *builderService) buildStaticBuildpackPack(
	ctx context.Context,
	st *state,
	bc *domain.BuildConfigStaticBuildpack,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	buildDir := filepath.Join(st.repositoryTempDir, contextDir)
	path, err := s.buildpack.Pack(ctx, buildDir, st.Logger(), s.tmpDestImage(st.app, st.build))
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
	_ context.Context,
	st *state,
) error {
	imageName := s.config.PartialTmpImageName(st.app.ID)
	digest, err := s.registry.ManifestDigest(imageName, st.build.Commit)
	if err != nil {
		return err
	}
	return s.registry.DeleteManifest(imageName, digest)
}
