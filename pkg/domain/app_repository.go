package domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	giturls "github.com/whilp/git-urls"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type Repository struct {
	ID       string
	Name     string
	URL      string
	Auth     optional.Of[RepositoryAuth]
	OwnerIDs []string
}

func NewRepository(name, url string, auth optional.Of[RepositoryAuth], ownerIDs []string) *Repository {
	return &Repository{
		ID:       NewID(),
		Name:     name,
		URL:      url,
		Auth:     auth,
		OwnerIDs: ownerIDs,
	}
}

func (r *Repository) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	ep, err := transport.NewEndpoint(r.URL)
	if err != nil {
		return errors.Wrap(err, "invalid url")
	}
	if !r.Auth.Valid {
		// URL is in http(s) format
		if ep.Protocol != "http" && ep.Protocol != "https" {
			return errors.New("url has to be http(s) protocol when auth is none")
		}
	} else if r.Auth.V.Method == RepositoryAuthMethodBasic {
		// URL is in https format
		if ep.Protocol != "https" {
			return errors.New("url has to be https protocol when auth is basic")
		}
	} else if r.Auth.V.Method == RepositoryAuthMethodSSH {
		// URL is in ssh format
		if ep.Protocol != "ssh" {
			return errors.New("url has to be ssh protocol when auth is ssh")
		}
	}
	if len(r.OwnerIDs) == 0 {
		return errors.New("owner_ids cannot be empty")
	}
	return nil
}

func (r *Repository) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(r.OwnerIDs, user.ID)
}

func (r *Repository) CanCreateApp(user *User) bool {
	// Only check for repository owner if repository is private;
	// allow everyone to create application if repository is public
	return r.IsOwner(user) || !r.Auth.Valid
}

// HTMLURL returns human-readable HTML page URL.
// Since there is no 'official' mapping from git-readable URL to human-readable HTML URL,
// the conversion is done heuristically.
func (r *Repository) HTMLURL() string {
	if !r.Auth.Valid {
		return r.URL // Expect a human-readable page
	}
	switch r.Auth.V.Method {
	case RepositoryAuthMethodBasic:
		return r.URL // Expect a human-readable page
	case RepositoryAuthMethodSSH:
		u, err := giturls.Parse(r.URL)
		if err != nil {
			return r.URL // Fallback
		}
		path := u.Path
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return "https://" + u.Hostname() + path + u.RawQuery
	default:
		return r.URL
	}
}

func (r *Repository) ResolveRefs(ctx context.Context, fallbackKey *ssh.PublicKeys) (refToCommit map[string]string, err error) {
	auth, err := GitAuthMethod(r, fallbackKey)
	if err != nil {
		return nil, err
	}
	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{r.URL},
	})
	refs, err := remote.ListContext(ctx, &git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to list remote refs at %v", r.URL))
	}

	refToCommit = make(map[string]string, 2*len(refs))
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

type RepositoryAuthMethod int

const (
	RepositoryAuthMethodBasic RepositoryAuthMethod = iota
	RepositoryAuthMethodSSH
)

type RepositoryAuth struct {
	Method   RepositoryAuthMethod
	Username string
	Password string
	SSHKey   string
}

func (r *RepositoryAuth) Validate() error {
	switch r.Method {
	case RepositoryAuthMethodBasic:
		if r.Username == "" {
			return errors.New("username cannot be empty")
		}
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
	case RepositoryAuthMethodSSH:
		if r.SSHKey != "" {
			_, err := ssh.NewPublicKeys("", []byte(r.SSHKey), "")
			if err != nil {
				return errors.Wrap(err, "invalid ssh private key")
			}
		}
	}
	return nil
}
