package domain

import (
	"github.com/friendsofgo/errors"
	"golang.org/x/crypto/ssh"
)

type AvatarBaseURL string

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

func (u *User) AvatarURL(baseURL AvatarBaseURL) string {
	return string(baseURL) + u.Name
}

type UserKey struct {
	ID        string
	UserID    string
	PublicKey string
}

func NewUserKey(userID string, publicKey string) (*UserKey, error) {
	key := &UserKey{
		ID:        NewID(),
		UserID:    userID,
		PublicKey: publicKey,
	}
	if err := key.Validate(); err != nil {
		return nil, err
	}
	return key, nil
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
