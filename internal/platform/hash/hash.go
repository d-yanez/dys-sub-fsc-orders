package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}
