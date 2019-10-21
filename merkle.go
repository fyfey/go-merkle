package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash generates a sha256 hash of given data
func Hash(in []byte) string {
	h := sha256.New()
	h.Write(in)
	return hex.EncodeToString(h.Sum(nil))
}
