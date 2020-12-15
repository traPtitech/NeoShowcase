package idgen

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// New 22文字以内のランダムな文字列を生成
func New() string {
	var b [16]byte
	_, _ = io.ReadFull(rand.Reader, b[:])
	return base64.RawURLEncoding.EncodeToString(b[:])
}
