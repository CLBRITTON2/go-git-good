package objects

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type Blob struct {
	Hash string
	Data []byte
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

// Create blob storage format for writing
func (blob *Blob) Serialize() []byte {
	header := fmt.Sprintf("blob %d\x00", len(blob.Data))
	data := append([]byte(header), blob.Data...)
	return data
}

// Create blob format from storage format
func DeserializeBlob(rawData []byte, hash string) (*Blob, error) {
	nullIndex := bytes.IndexByte(rawData, byte('\x00'))
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid blob format: no null byte found")
	}

	header := string(rawData[:nullIndex])
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "blob" {
		return nil, fmt.Errorf("invalid blob header format: %s", header)
	}

	data := rawData[nullIndex+1:]
	blob := &Blob{
		Hash: hash,
		Data: data,
	}

	return blob, nil
}

func CalculateHash(encodedData []byte) string {
	hasher := sha1.New()
	hasher.Write(encodedData)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
