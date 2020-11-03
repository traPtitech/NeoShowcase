package idgen

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func New() string {
	var b [16]byte
	_, _ = io.ReadFull(rand.Reader, b[:])
	return base64.RawURLEncoding.EncodeToString(b[:])
}
