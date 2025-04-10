package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	"github.com/traPtitech/neoshowcase/pkg/domain"
)

var _ domain.GitService = (*service)(nil)

type service struct {
	fallbackKey *ssh.PublicKeys
}

func NewService(fallbackKey *ssh.PublicKeys) domain.GitService {
	return &service{fallbackKey: fallbackKey}
}

func (s *service) gitAuth(repo *domain.Repository) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	if repo.Auth.Valid {
		switch repo.Auth.V.Method {
		case domain.RepositoryAuthMethodBasic:
			auth = &http.BasicAuth{
				Username: repo.Auth.V.Username,
				Password: repo.Auth.V.Password,
			}
		case domain.RepositoryAuthMethodSSH:
			if repo.Auth.V.SSHKey != "" {
				keys, err := domain.IntoPublicKey(domain.PrivateKey(repo.Auth.V.SSHKey))
				if err != nil {
					return nil, err
				}
				auth = keys
			} else {
				auth = s.fallbackKey
			}
		}
	}
	return auth, nil
}

func (s *service) ResolveRefs(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
	auth, err := s.gitAuth(repo)
	if err != nil {
		return nil, err
	}

	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repo.URL},
	})
	refs, err := remote.ListContext(ctx, &git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("list remote refs at %v", repo.URL))
	}

	refToCommit := make(map[string]string, 2*len(refs))
	for _, ref := range refs {
		if ref.Type() == plumbing.HashReference {
			refToCommit[ref.Name().String()] = ref.Hash().String()
			refToCommit[ref.Name().Short()] = ref.Hash().String()
		}
	}
	for _, ref := range refs {
		if ref.Type() == plumbing.SymbolicReference {
			commit, ok := refToCommit[ref.Target().String()]
			if ok {
				refToCommit[ref.Name().String()] = commit
			}
		}
	}

	return refToCommit, nil
}

func (s *service) CloneRepository(ctx context.Context, dir string, repo *domain.Repository, commitHash string) error {
	auth, err := s.gitAuth(repo)
	if err != nil {
		return err
	}

	localRepo, remote, err := initializeGitRepo(dir, repo.URL, false)
	if err != nil {
		return err
	}
	targetRef := plumbing.NewRemoteReferenceName("origin", "target")
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", commitHash, targetRef))},
		Depth:      1,
		Auth:       auth,
	})
	if err != nil {
		return errors.Wrap(err, "fetch commit")
	}

	wt, err := localRepo.Worktree()
	if err != nil {
		return errors.Wrap(err, "get worktree")
	}
	err = wt.Checkout(&git.CheckoutOptions{Branch: targetRef})
	if err != nil {
		return errors.Wrap(err, "checkout")
	}

	err = updateSubModules(wt, auth)
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(dir, ".git"))
	if err != nil {
		return err
	}

	return nil
}

func updateSubModules(wt *git.Worktree, auth transport.AuthMethod) error {
	sm, err := wt.Submodules()
	if err != nil {
		return errors.Wrap(err, "get submodules")
	}
	// Try with auth first, then try without auth
	err = sm.Update(&git.SubmoduleUpdateOptions{
		Init:              true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
		Depth:             1,
	})
	if err != nil {
		err = sm.Update(&git.SubmoduleUpdateOptions{
			Init:              true,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Depth:             1,
		})
		if err != nil {
			return errors.Wrap(err, "update submodules")
		}
	}
	return nil
}

func (s *service) CreateBareRepository(dir string, repo *domain.Repository) (domain.GitRepository, error) {
	localRepo, remote, err := initializeGitRepo(dir, repo.URL, true)
	if err != nil {
		return nil, err
	}
	auth, err := s.gitAuth(repo)
	if err != nil {
		return nil, err
	}
	return &repository{repo: localRepo, remote: remote, auth: auth}, nil
}

func initializeGitRepo(dir, remoteURL string, isBare bool) (*git.Repository, *git.Remote, error) {
	localRepo, err := git.PlainInit(dir, isBare)
	if err != nil {
		return nil, nil, errors.Wrap(err, "init repository")
	}
	remote, err := localRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "create remote")
	}
	return localRepo, remote, nil
}

var _ domain.GitRepository = (*repository)(nil)

type repository struct {
	repo   *git.Repository
	remote *git.Remote
	auth   transport.AuthMethod
}

func (r *repository) Fetch(ctx context.Context, hashes []string) error {
	refSpecs := lo.Map(hashes, func(hash string, i int) config.RefSpec {
		targetRef := plumbing.NewRemoteReferenceName("origin", fmt.Sprintf("target-%d", i))
		return config.RefSpec(fmt.Sprintf("+%s:%s", hash, targetRef))
	})
	err := r.remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   refSpecs,
		Depth:      1,
		Auth:       r.auth,
	})
	if err != nil {
		return errors.Wrap(err, "fetch commits")
	}

	return nil
}

func (r *repository) GetCommit(hash string) (*domain.RepositoryCommit, error) {
	commit, err := r.repo.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "get commit")
	}
	return toRepositoryCommit(commit), nil
}

func toRepositoryCommit(c *object.Commit) *domain.RepositoryCommit {
	return &domain.RepositoryCommit{
		Hash: c.Hash.String(),
		Author: domain.CommitAuthorSignature{
			Name:  c.Author.Name,
			Email: c.Author.Email,
			Date:  c.Author.When,
		},
		Committer: domain.CommitAuthorSignature{
			Name:  c.Committer.Name,
			Email: c.Committer.Email,
			Date:  c.Committer.When,
		},
		Message: c.Message,
		Error:   false,
	}
}
