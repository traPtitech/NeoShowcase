package random

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
)

func SecureGenerate(charset []rune, length int) string {
	max := big.NewInt(int64(len(charset)))
	var b strings.Builder
	b.Grow(length)
	for i := 0; i < length; i++ {
		charIdx, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		b.WriteRune(charset[charIdx.Int64()])
	}
	return b.String()
}

const (
	lowerCharSet    = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolCharSet   = "!@#$%&*"
	numberSet       = "0123456789"
	passwordCharSet = lowerCharSet + upperCharSet + symbolCharSet + numberSet
)

var (
	passwordRunes = []rune(passwordCharSet)
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

func SecureGeneratePassword(length int) string {
	return SecureGenerate(passwordRunes, length)
}
