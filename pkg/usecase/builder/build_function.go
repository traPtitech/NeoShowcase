package builder

import (
	"context"
	"path/filepath"

	buildkit "github.com/moby/buildkit/client"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

// buildFunctionBuildpack builds a function using buildpack.
func (s *ServiceImpl) buildFunctionBuildpack(
	ctx context.Context,
	st *state,
	bc *domain.BuildConfigFunctionBuildpack,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	buildDir := filepath.Join(st.repositoryTempDir, contextDir)
	path, err := s.buildpack.Pack(ctx, buildDir, s.tmpDestImage(st.app, st.build), s.imageConfig, st.appEnv(), st.Logger())
	if err != nil {
		return err
	}
	st.staticDest = filepath.Join(path, bc.ArtifactPath)
	return nil
}

// buildFunctionCmd builds a function using command.
func (s *ServiceImpl) buildFunctionCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigFunctionCmd,
) error {
	// Similar to buildStaticCmd but for function
	return s.buildStaticCmd(ctx, st, ch, &domain.BuildConfigStaticCmd{
		StaticConfig: domain.StaticConfig{
			ArtifactPath: bc.ArtifactPath,
			SPA:          false,
		},
		BaseImage: bc.BaseImage,
		BuildCmd:  bc.BuildCmd,
	})
}

// buildFunctionDockerfile builds a function using Dockerfile.
func (s *ServiceImpl) buildFunctionDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigFunctionDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	st.staticDest = bc.ArtifactPath
	return s.solveDockerfile(
		ctx,
		s.tmpDestImage(st.app, st.build),
		filepath.Join(st.repositoryTempDir, contextDir),
		filepath.Join(st.repositoryTempDir, contextDir),
		bc.DockerfileName,
		st.appEnv(),
		ch,
	)
}

// extractFunctionArtifact extracts the single JS file from the built image.
func (s *ServiceImpl) extractFunctionArtifact(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	// For now, use the same extraction logic as static builds
	// This extracts the artifact from the temporary image
	return s.buildStaticExtract(ctx, st, ch)
}
