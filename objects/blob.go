package objects

import (
	"fmt"
	"os"

	"github.com/CLBRITTON2/go-git-good/common"
)

type Blob struct {
	Hash common.Hash
	Data []byte
}

func CreateBlobFromFile(fileToBlob string) (*Blob, error) {
	data, err := os.ReadFile(fileToBlob)
	if err != nil {
		return nil, fmt.Errorf("error reading file for blob: %v", err)
	}
	newBlob := &Blob{
		Data: data,
	}
	serializedBlobData := newBlob.Serialize()
	hash, err := common.HashObject(serializedBlobData)
	if err != nil {
		return nil, err
	}
	if hash.Empty() {
		return nil, fmt.Errorf("error creating blob from file: HashObject returned an empty hash")
	}
	newBlob.Hash = hash
	return newBlob, nil
}

// Create blob storage format for writing
func (blob *Blob) Serialize() []byte {
	header := fmt.Sprintf("blob %d\x00", len(blob.Data))
	data := append([]byte(header), blob.Data...)
	return data
}
