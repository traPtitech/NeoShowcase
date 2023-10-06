package builder

import (
	"context"
	"path/filepath"

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
	st.staticDest = filepath.Join(path, bc.ArtifactPath)
	return nil
}
