package domain

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/friendsofgo/errors"
	ssh2 "golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func NewPublicKey(pemBytes []byte) (*ssh.PublicKeys, error) {
	return ssh.NewPublicKeys("git", pemBytes, "")
}

func Base64EncodedPublicKey(pubKey ssh2.PublicKey) string {
	return pubKey.Type() + " " + base64.StdEncoding.EncodeToString(pubKey.Marshal())
}

func EncodePrivateKeyPem(privKey ed25519.PrivateKey) (string, error) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", errors.Wrap(err, "encoding private key")
	}
	var res bytes.Buffer
	privateKeyPEM := &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}
	err = pem.Encode(&res, privateKeyPEM)
	if err != nil {
		return "", errors.Wrap(err, "encoding pem")
	}
	return res.String(), nil
}

func GitAuthMethod(repo *Repository, fallbackKey *ssh.PublicKeys) (transport.AuthMethod, error) {
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
				keys, err := NewPublicKey([]byte(repo.Auth.V.SSHKey))
				if err != nil {
					return nil, err
				}
				auth = keys
			} else {
				auth = fallbackKey
			}
		}
	}
	return auth, nil
}
