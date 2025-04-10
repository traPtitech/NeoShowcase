package builder

import (
	"context"
)

func (s *ServiceImpl) cloneRepository(ctx context.Context, st *state) error {
	return s.gitsvc.CloneRepository(ctx, st.repositoryTempDir, st.repo, st.build.Commit)
}
