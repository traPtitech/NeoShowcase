package domain

import (
	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"

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