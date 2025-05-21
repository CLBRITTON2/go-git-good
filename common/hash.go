package common

import (
	"crypto/sha1"
	"encoding/hex"
)

type Hash [20]byte

func HashObject(data []byte) (Hash, error) {
	hasher := sha1.New()
	hasher.Write(data)
	sha1 := hasher.Sum(nil)
	var hash Hash
	copy(hash[:], sha1)
	return hash, nil
}

func (hash Hash) String() string {
	return hex.EncodeToString(hash[:])
}

func (hash Hash) Empty() bool {
	return hash == Hash{}
}
