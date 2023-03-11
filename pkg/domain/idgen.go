package domain

import (
	"github.com/traPtitech/neoshowcase/pkg/util/random"
)

const IDLength = 22

// NewID 22文字のランダムな文字列を生成
func NewID() string {
	return random.SecureGenerateHex(IDLength)
}
