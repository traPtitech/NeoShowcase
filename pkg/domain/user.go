package domain

import (
	"github.com/friendsofgo/errors"
	"golang.org/x/crypto/ssh"
)

func NewUser(name string) *User {
	return &User{
		ID:    NewID(),
		Name:  name,
		Admin: false,
	}
}

type User struct {
	ID    string
	Name  string
	Admin bool
}

type UserKey struct {
	ID        string
	UserID    string
	PublicKey string
}

func (u *UserKey) MarshalKey() []byte {
	out, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(u.PublicKey))
	return out.Marshal()
}

func (u *UserKey) Validate() error {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(u.PublicKey))
	if err != nil {
		return errors.Wrap(err, "invalid public key format")
	}
	return nil
}
