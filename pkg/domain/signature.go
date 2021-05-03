package domain

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
)

func VerifySignature(h func() hash.Hash, message, key []byte, signature string) bool {
	decoded, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return validate(h, message, key, decoded)
}

func validate(h func() hash.Hash, message, key, signature []byte) bool {
	mac := hmac.New(h, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	return hmac.Equal(signature, sum)
}
