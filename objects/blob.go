package objects

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type Blob struct {
	Hash string
	Data []byte
}

// Create blob storage format for writing
func (blob *Blob) Serialize() []byte {
	header := fmt.Sprintf("blob %d\x00", len(blob.Data))
	data := append([]byte(header), blob.Data...)
	return data
}

func CreateBlobFromFile(fileToBlob string) (*Blob, error) {
	data, err := os.ReadFile(fileToBlob)
	if err != nil {
		return nil, fmt.Errorf("error reading file for blob: %v", fileToBlob)
	}
	newBlob := &Blob{
		Data: data,
	}

	encodedData := newBlob.Serialize()
	newBlob.Hash = CalculateHash(encodedData)
	return newBlob, nil
}

func CalculateHash(encodedData []byte) string {
	hasher := sha1.New()
	hasher.Write(encodedData)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
