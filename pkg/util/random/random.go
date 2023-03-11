package random

import (
	"crypto/rand"
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
	hex             = "0123456789abcdef"
	lowerCharSet    = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolCharSet   = "!@#$%&*"
	numberSet       = "0123456789"
	passwordCharSet = lowerCharSet + upperCharSet + symbolCharSet + numberSet
)

var (
	hexRunes      = []rune(hex)
	passwordRunes = []rune(passwordCharSet)
)

func SecureGenerateHex(length int) string {
	return SecureGenerate(hexRunes, length)
}

func SecureGeneratePassword(length int) string {
	return SecureGenerate(passwordRunes, length)
}
