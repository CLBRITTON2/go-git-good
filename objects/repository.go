package objects

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	WorkTree     string
	GitDirectory string
}

func CreateRepository(path string) (*Repository, error) {
	if path == "" {
		path = "."
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("could not resolve absolute path: %w", err)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %v", err)
	}

	repository := Repository{
		WorkTree:     path,
		GitDirectory: filepath.Join(path, ".gitgood"),
	}

	if info, err := os.Stat(repository.GitDirectory); err == nil && info.IsDir() {
		return nil, fmt.Errorf("git repository already exists at %v", path)
	}

	fmt.Printf("Creating an empty git repository at: %v\n", path)
	err = os.MkdirAll(repository.GitDirectory, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating .git directory: %v", err)
	}

	// Create repository subdirectories
	repositoryDirectories := []string{"objects", "refs/heads", "refs/tags"}
	for _, directory := range repositoryDirectories {
		fullPath := filepath.Join(repository.GitDirectory, directory)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %v", err)
		}
	}

	defaultConfigContents := []byte("[core]\n repositoryformatversion = 0\n filemode = true\n bare = false\n")
	err = os.WriteFile(filepath.Join(repository.GitDirectory, "config"), defaultConfigContents, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing default config file: %v", err)
	}

	defaultHeadContents := []byte("ref: refs/heads/master\n")
	err = os.WriteFile(filepath.Join(repository.GitDirectory, "HEAD"), defaultHeadContents, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing HEAD file: %v", err)
	}

	defaultDescriptionContents := []byte("Unnamed repository: edit 'description' file to name the repository\n")
	err = os.WriteFile(filepath.Join(repository.GitDirectory, "description"), defaultDescriptionContents, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing description file: %v", err)
	}

	return &repository, nil
}

func FindRepository(path string) (*Repository, error) {
	if path == "" {
		path = "."
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("could not resolve absolute path: %w", err)
	}

	gitDirectory := filepath.Join(path, ".gitgood")
	if info, err := os.Stat(gitDirectory); err == nil && info.IsDir() {
		repository := Repository{
			WorkTree:     path,
			GitDirectory: gitDirectory,
		}
		return &repository, nil
	}

	parentDirectory := filepath.Dir(path)

	if parentDirectory == path {
		return nil, fmt.Errorf("no git repository found")
	}

	return FindRepository(parentDirectory)
}

func (repository *Repository) WriteObject(objectHash string, serializedData []byte) error {
	// The subdirectory to write the object to is the first 2 characters of the SHA-1
	// The reamining 38 characters are the filename
	directory := objectHash[0:2]
	path := filepath.Join(repository.GitDirectory, "objects")
	objectDirectory := filepath.Join(path, directory)
	err := os.MkdirAll(objectDirectory, 0755)
	if err != nil {
		return fmt.Errorf("error making write object directory: %v", err)
	}

	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(serializedData)
	if err != nil {
		return fmt.Errorf("error compressing object data: %v", err)
	}
	writer.Close()

	hashFileName := objectHash[2:]
	objectFilePath := filepath.Join(objectDirectory, hashFileName)

	err = os.WriteFile(objectFilePath, buffer.Bytes(), 0664)
	if err != nil {
		return fmt.Errorf("error writing compressed data to file: %v", err)
	}
	return nil
}

func (repository *Repository) ReadObject(objectHash string) (string, error) {
	directory := objectHash[0:2]
	path := filepath.Join(repository.GitDirectory, "objects")
	objectDirectory := filepath.Join(path, directory)

	hashFileName := objectHash[2:]
	objectFilePath := filepath.Join(objectDirectory, hashFileName)
	compressedData, err := os.ReadFile(objectFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading compressed object file: %v", err)
	}

	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", fmt.Errorf("error creating zlib reader %v", err)
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return "", fmt.Errorf("error decompressing object data: %v", err)
	}

	nullIndex := bytes.IndexByte(decompressedData.Bytes(), byte('\x00'))
	if nullIndex == -1 {
		return "", fmt.Errorf("invalid object format: no null byte found")
	}

	header := string(decompressedData.Bytes()[:nullIndex])
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid object header format expected <type> <data length> got: %s", header)
	}
	objectType := parts[0]

	switch objectType {
	case "blob":
		blob, err := DeserializeBlob(decompressedData.Bytes(), objectHash)
		if err != nil {
			return "", fmt.Errorf("error deserializing blob: %v", err)
		}
		return string(blob.Data), nil
	default:
		return "", fmt.Errorf("unknown object type: %s", objectType)
	}
}
