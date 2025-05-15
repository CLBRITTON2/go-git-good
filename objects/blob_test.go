package objects

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CLBRITTON2/go-git-good/common"
)

func createTestFile(t *testing.T, tempDir, content string) string {
	t.Helper()
	filePath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file for create blob from file test: %v", err)
	}
	return filePath
}

func TestCreateBlobFromFileValid(t *testing.T) {
	tempDir := t.TempDir()
	filePath := createTestFile(t, tempDir, "test")
	blob, err := CreateBlobFromFile(filePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(blob.Data) != "test" {
		t.Errorf("expected Data %q, got %q", "test", string(blob.Data))
	}
	expectedHash, err := common.HashObject("blob", []byte("test"))
	if err != nil {
		t.Fatalf("unexpected error calculating hash: %v", err)
	}
	if blob.Hash != expectedHash {
		t.Errorf("expected Hash %q, got %q", expectedHash.String(), blob.Hash.String())
	}
}

func TestCreateBlobFromFileEmpty(t *testing.T) {
	tempDir := t.TempDir()
	filePath := createTestFile(t, tempDir, "")
	blob, err := CreateBlobFromFile(filePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(blob.Data) != 0 {
		t.Errorf("expected empty Data, got %q", string(blob.Data))
	}
	expectedHash, err := common.HashObject("blob", []byte{})
	if err != nil {
		t.Fatalf("unexpected error calculating hash: %v", err)
	}
	if blob.Hash != expectedHash {
		t.Errorf("expected Hash %q, got %q", expectedHash.String(), blob.Hash.String())
	}
}

func TestCreateBlobFromFileNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "nonexistent.txt")
	blob, err := CreateBlobFromFile(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if blob != nil {
		t.Errorf("expected nil Blob, got %v", blob)
	}
	if !strings.Contains(err.Error(), "error reading file for blob") {
		t.Errorf("expected error containing 'error reading file for blob', got %q", err.Error())
	}
}
