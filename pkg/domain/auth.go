package domain

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/friendsofgo/errors"
	ssh2 "golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type PrivateKey []byte

func IntoPublicKey(pk PrivateKey) (*ssh.PublicKeys, error) {
	return ssh.NewPublicKeys("git", pk, "")
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
