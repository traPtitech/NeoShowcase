package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

const validSSHKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMwwAAAJCX+efxl/nn
8QAAAAtzc2gtZWQyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMww
AAAEA+FzwWKIYduEDOqkEOZ2wmxZWPc2wpZeWj+J8e3Q6x0QLWIALnhPWigJA3njldxcM9
1PEkQQN231Ow7czCg8zDAAAADG1vdG9AbW90by13cwE=
-----END OPENSSH PRIVATE KEY-----`

func TestRepository_Validate(t *testing.T) {
	tests := []struct {
		name    string
		repo    Repository
		wantErr bool
	}{
		{
			name: "valid auth none (http)",
			repo: Repository{
				Name:     "test",
				URL:      "http://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			wantErr: false,
		},
		{
			name: "valid auth none (https)",
			repo: Repository{
				Name:     "test",
				URL:      "https://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			wantErr: false,
		},
		{
			name: "valid auth basic",
			repo: Repository{
				Name: "test",
				URL:  "https://github.com/traPtitech/NeoShowcase",
				Auth: optional.From(RepositoryAuth{
					Method:   RepositoryAuthMethodBasic,
					Username: "username",
					Password: "password",
				}),
				OwnerIDs: []string{"abc"},
			},
			wantErr: false,
		},
		{
			name: "valid auth ssh",
			repo: Repository{
				Name: "test",
				URL:  "git@github.com:traPtitech/NeoShowcase.git",
				Auth: optional.From(RepositoryAuth{
					Method: RepositoryAuthMethodSSH,
					SSHKey: validSSHKey,
				}),
				OwnerIDs: []string{"abc"},
			},
			wantErr: false,
		},
		{
			name: "invalid name",
			repo: Repository{
				Name:     "",
				URL:      "http://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid url",
			repo: Repository{
				Name:     "test",
				URL:      "ttp://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid scheme (auth none)",
			repo: Repository{
				Name:     "test",
				URL:      "git@github.com:traPtitech/NeoShowcase.git",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid scheme (auth basic)",
			repo: Repository{
				Name: "test",
				URL:  "http://github.com/traPtitech/NeoShowcase",
				Auth: optional.From(RepositoryAuth{
					Method:   RepositoryAuthMethodBasic,
					Username: "username",
					Password: "password",
				}),
				OwnerIDs: []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid scheme (auth ssh)",
			repo: Repository{
				Name: "test",
				URL:  "https://github.com/traPtitech/NeoShowcase",
				Auth: optional.From(RepositoryAuth{
					Method: RepositoryAuthMethodSSH,
					SSHKey: validSSHKey,
				}),
				OwnerIDs: []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid owners",
			repo: Repository{
				Name:     "test",
				URL:      "http://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepositoryAuth_Validate(t *testing.T) {
	tests := []struct {
		name    string
		auth    RepositoryAuth
		wantErr bool
	}{
		{
			name: "valid basic auth",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodBasic,
				Username: "root",
				Password: "password",
				SSHKey:   "",
			},
			wantErr: false,
		},
		{
			name: "invalid username",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodBasic,
				Username: "",
				Password: "password",
				SSHKey:   "",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodBasic,
				Username: "root",
				Password: "",
				SSHKey:   "",
			},
			wantErr: true,
		},
		{
			name: "valid ssh auth",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey:   validSSHKey,
			},
			wantErr: false,
		},
		{
			name: "valid ssh auth (uses default system key)",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey:   "",
			},
			wantErr: false,
		},
		{
			name: "invalid ssh private key",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey: `-----BEGIN OPENSSH PRIVATE KEY------
-----END OPENSSH PRIVATE KEY-----`,
			},
			wantErr: true,
		},
		{
			name: "invalid ssh auth (public key)",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey:   `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALWIALnhPWigJA3njldxcM91PEkQQN231Ow7czCg8zD`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auth.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
