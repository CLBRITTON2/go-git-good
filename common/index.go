package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type Index struct {
	NumberOfEntries uint32
	Entries         []*IndexEntry
}

type IndexEntry struct {
	ModifiedTime time.Time
	Hash         Hash
	FileSize     uint32
	FileMode     uint32
	EntryPath    string
}

func (index *Index) Exists(repository *Repository) (bool, error) {
	indexPath := filepath.Join(repository.GitDirectory, "index")
	_, err := os.ReadFile(indexPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetIndex(repository *Repository) (*Index, error) {
	indexPath := filepath.Join(repository.GitDirectory, "index")
	index, err := ReadIndex(indexPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Index{
				NumberOfEntries: 0,
				Entries:         []*IndexEntry{},
			}, nil
		}
		return nil, err
	}
	return index, nil
}

func ReadIndex(filePath string) (*Index, error) {
	indexFileData, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, fmt.Errorf("%v", err)
	}

	if string(indexFileData[:4]) != "DIRC" {
		return nil, fmt.Errorf("error reading index header: invalid signature")
	}
	indexVersion := binary.BigEndian.Uint32(indexFileData[4:8])
	if indexVersion != 2 {
		return nil, fmt.Errorf("error reading index version number: expected 2 got %v", indexVersion)
	}
	indexEntryCount := binary.BigEndian.Uint32(indexFileData[8:12])

	index := &Index{
		NumberOfEntries: indexEntryCount,
		Entries:         make([]*IndexEntry, 0, indexEntryCount),
	}

	// Entries will live 12 bytes past the header
	buffer := bytes.NewReader(indexFileData[12:])
	for i := uint32(0); i < indexEntryCount; i++ {
		// 36 bytes are required minimum (mod time 8, file mode 4, file size 4, hash 20)
		// If we only have 36 bytes we're missing something
		// Entry path is variable
		if buffer.Len() <= 36 {
			return nil, fmt.Errorf("error reading index entry: content less than 36 bytes for required fields")
		}
		var modifiedTimeSeconds, modifiedTimeNano uint32
		binary.Read(buffer, binary.BigEndian, &modifiedTimeSeconds)
		binary.Read(buffer, binary.BigEndian, &modifiedTimeNano)
		modifiedTime := time.Unix(int64(modifiedTimeSeconds), int64(modifiedTimeNano))

		var fileMode uint32
		binary.Read(buffer, binary.BigEndian, &fileMode)

		var fileSize uint32
		binary.Read(buffer, binary.BigEndian, &fileSize)

		var hash [20]byte
		binary.Read(buffer, binary.BigEndian, &hash)

		var pathBytes []byte
		for {
			pathBuffer, err := buffer.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("error reading index entry path: %v", err)
			}
			// Entry paths are null terminated
			if pathBuffer == 0 {
				break
			}
			pathBytes = append(pathBytes, pathBuffer)
		}
		indexEntryPath := string(pathBytes)

		indexEntry := &IndexEntry{
			ModifiedTime: modifiedTime,
			FileMode:     fileMode,
			FileSize:     fileSize,
			Hash:         hash,
			EntryPath:    indexEntryPath,
		}

		index.Entries = append(index.Entries, indexEntry)
	}
	return index, nil
}

func WriteIndex(repository *Repository, index *Index) error {
	// DIRC = dircache in normal Git index files
	// This is a dumbed down version but may as well keep the header consistent for now
	header := [4]byte{'D', 'I', 'R', 'C'}
	version := uint32(2)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, header)
	binary.Write(&buffer, binary.BigEndian, version)
	binary.Write(&buffer, binary.BigEndian, index.NumberOfEntries)

	// Mimic git index format as close as possible excluding some metadata ie dev/ino/uid etc
	for _, entry := range index.Entries {
		modifiedTimeSeconds := uint32(entry.ModifiedTime.Unix())
		modifiedTimeNano := uint32(entry.ModifiedTime.Nanosecond())
		binary.Write(&buffer, binary.BigEndian, modifiedTimeSeconds)
		binary.Write(&buffer, binary.BigEndian, modifiedTimeNano)
		binary.Write(&buffer, binary.BigEndian, entry.FileMode)
		binary.Write(&buffer, binary.BigEndian, entry.FileSize)
		binary.Write(&buffer, binary.BigEndian, entry.Hash)
		buffer.WriteString(entry.EntryPath)
		buffer.WriteByte(0)
	}

	indexPath := filepath.Join(repository.GitDirectory, "index")
	err := os.WriteFile(indexPath, buffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing data to index file: %v", err)
	}
	return nil
}

func (index *Index) AddEntry(entry *IndexEntry) {
	for i, existingEntry := range index.Entries {
		if existingEntry.EntryPath == entry.EntryPath {
			index.Entries[i] = entry
			return
		}
	}

	index.Entries = append(index.Entries, entry)
	index.NumberOfEntries = uint32(len(index.Entries))
}

func (index *Index) RemoveEntry(entryPath string) {
	for i, entry := range index.Entries {
		if entry.EntryPath == entryPath {
			index.Entries = slices.Delete(index.Entries, i, i+1)
			index.NumberOfEntries = uint32(len(index.Entries))
			return
		}
	}
}
