package random

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

func SecureGenerateHex(length int) string {
	byteLength := (length + 1) / 2
	b := make([]byte, byteLength)
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

// SecureGeneratePassword generates a secure random password of the specified length.
// The password is base64 encoded to ensure it is URL-safe and can include a wide range
func SecureGeneratePassword(length int) string {
	b := make([]byte, length*3/4+1) // to ensure enough bytes for base64 encoding
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}
