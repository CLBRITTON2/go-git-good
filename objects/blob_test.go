package objects

import (
	"os"
	"path/filepath"
	"testing"
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
	if err != nil {
		t.Fatalf("unexpected error calculating hash: %v", err)
	}

	// Matches git hash-object for a file named test.text with content "test" no new line at the end
	expectedHash := "30d74d258442c7c65512eafab474568dd706c430"
	if blob.Hash.String() != expectedHash {
		t.Errorf("expected Hash %q, got %q", expectedHash, blob.Hash.String())
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

	expectedHash := "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391"
	if blob.Hash.String() != expectedHash {
		t.Errorf("expected Hash %q, got %q", expectedHash, blob.Hash.String())
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
}
