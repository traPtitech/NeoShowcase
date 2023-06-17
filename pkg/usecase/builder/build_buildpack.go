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
	return s.buildpack.Pack(ctx, buildDir, st.Logger(), s.destImage(st.app, st.build))
}
