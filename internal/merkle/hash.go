package merkle

import "crypto/sha256"

type Hasher interface {
	Hash(data []byte) []byte
}

type SHA256Hasher struct{}

func (s SHA256Hasher) Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

type DoubleSHA256Hasher struct{}

func (d DoubleSHA256Hasher) Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])
	return hash[:]
}
