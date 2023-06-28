package hash

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/zeebo/xxh3"
)

// XXH3Hex returns 64-bit xxh3 hash in hexadecimal form (16 chars).
func XXH3Hex(b []byte) string {
	var res [8]byte
	binary.LittleEndian.PutUint64(res[:], xxh3.Hash(b))
	return hex.EncodeToString(res[:])
}
