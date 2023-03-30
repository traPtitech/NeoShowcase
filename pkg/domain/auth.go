package domain

import (
	"encoding/base64"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func Base64EncodedPublicKey(keys *ssh.PublicKeys) string {
	pubKey := keys.Signer.PublicKey()
	return pubKey.Type() + " " + base64.StdEncoding.EncodeToString(pubKey.Marshal())
}

func GitAuthMethod(repo *Repository, sshFallback *ssh.PublicKeys) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	if repo.Auth.Valid {
		switch repo.Auth.V.Method {
		case RepositoryAuthMethodBasic:
			auth = &http.BasicAuth{
				Username: repo.Auth.V.Username,
				Password: repo.Auth.V.Password,
			}
		case RepositoryAuthMethodSSH:
			if repo.Auth.V.SSHKey != "" {
				keys, err := ssh.NewPublicKeys("", []byte(repo.Auth.V.SSHKey), "")
				if err != nil {
					return nil, err
				}
				auth = keys
			} else {
				auth = sshFallback
			}
		}
	}
	return auth, nil
}
