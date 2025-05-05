package objects

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type Blob struct {
	Hash string
	Size int
	Data []byte
}

func CreateBlob(fileToBlob string) (*Blob, error) {
	data, err := os.ReadFile(fileToBlob)
	if err != nil {
		return nil, fmt.Errorf("error reading file for blob: %v", fileToBlob)
	}
	size := len(data)
	newBlob := &Blob{
		Size: size,
		Data: data,
	}

	newBlob.Hash = CalculateHash("blob", newBlob.Data)
	return newBlob, nil
}

func CalculateHash(objectType string, objectData []byte) string {
	header := fmt.Sprintf("%s %d\x00", objectType, len(objectData))
	data := append([]byte(header), objectData...)
	hasher := sha1.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
