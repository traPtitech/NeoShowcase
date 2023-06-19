package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (s *builderService) cloneRepository(ctx context.Context, st *state) error {
	repo, err := git.PlainInit(st.repositoryTempDir, false)
	if err != nil {
		return errors.Wrap(err, "failed to init repository")
	}
	auth, err := domain.GitAuthMethod(st.repo, s.pubKey)
	if err != nil {
		return err
	}
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{st.repo.URL},
	})
	if err != nil {
		return errors.Wrap(err, "failed to add remote")
	}
	targetRef := plumbing.NewRemoteReferenceName("origin", "target")
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", st.build.Commit, targetRef))},
		Depth:      1,
		Auth:       auth,
	})
	if err != nil {
		return errors.Wrap(err, "failed to clone repository")
	}
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}
	err = wt.Checkout(&git.CheckoutOptions{Branch: targetRef})
	if err != nil {
		return errors.Wrap(err, "failed to checkout")
	}
	sm, err := wt.Submodules()
	if err != nil {
		return errors.Wrap(err, "getting submodules")
	}

	// Try with auth first, then try without auth
	err = sm.Update(&git.SubmoduleUpdateOptions{
		Init:              true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})
	if err != nil {
		err = sm.Update(&git.SubmoduleUpdateOptions{
			Init:              true,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			return errors.Wrap(err, "updating submodules")
		}
	}

	// Delete .git directory before passing to the builder
	err = os.RemoveAll(filepath.Join(st.repositoryTempDir, ".git"))
	if err != nil {
		return errors.Wrap(err, "deleting .git directory")
	}
	return nil
}
