package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Hash [20]byte

var validObjectTypes = map[string]bool{
	"blob": true,
	"tree": true,
}

func HashObject(objectType string, data []byte) (Hash, error) {
	if !validObjectTypes[objectType] {
		return Hash{}, fmt.Errorf("invalid object type at HashObject: %s", objectType)
	}

	header := fmt.Sprintf("%s %d\x00", objectType, len(data))
	serializedData := append([]byte(header), data...)
	hasher := sha1.New()
	hasher.Write(serializedData)
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
