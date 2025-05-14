package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Hash [20]byte

func HashObject(objectType string, data []byte) Hash {
	header := fmt.Sprintf("%s %d\x00", objectType, len(data))
	serializedData := append([]byte(header), data...)
	hasher := sha1.New()
	hasher.Write(serializedData)
	sha1 := hasher.Sum(nil)
	var hash Hash
	copy(hash[:], sha1)
	return hash
}

func (hash Hash) String() string {
	return hex.EncodeToString(hash[:])
}
