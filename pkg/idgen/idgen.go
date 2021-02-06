package idgen

import (
	"math/rand"
)

// New 22文字のランダムな文字列を生成
func New() string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	str := make([]byte, 22)
	for i := 0; i < 22; i++ {
		str[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(str)
}
