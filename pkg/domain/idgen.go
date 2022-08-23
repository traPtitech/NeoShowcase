package domain

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

const IDLength = 22

// NewID 22文字のランダムな文字列を生成
func NewID() string {
	const bytesLength = IDLength / 2
	b := make([]byte, bytesLength)
	n, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if n != bytesLength {
		panic(fmt.Errorf("expected %d bytes, but got %d bytes", bytesLength, n))
	}
	return hex.EncodeToString(b)
}
