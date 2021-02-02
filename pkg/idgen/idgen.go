package idgen

import (
	"github.com/volatiletech/randomize"
	"math/rand"
)

// New 22文字以内のランダムな文字列を生成
func New() string {
	return randomize.Str(rand.Int63, 22)
}
