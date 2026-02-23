package builder

import "context"

func (s *ServiceImpl) buildRegistryCleanup(
	ctx context.Context,
	st *state,
) error {
	return s.regclient.DeleteImage(ctx, s.imageConfig.TmpImageName(st.app.ID), s.imageTag(st.build))
}
