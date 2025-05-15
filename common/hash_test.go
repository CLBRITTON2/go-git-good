package common

import "testing"

func TestHashObjectSimple(t *testing.T) {
	// Data from echo -n "test" | git hash-object --stdin
	objectType := "blob"
	data := []byte("test")
	expected := "30d74d258442c7c65512eafab474568dd706c430"

	hashBytes, err := HashObject(objectType, data)
	if err != nil {
		t.Fatalf("HashObject(%q, %q) unexpected error: %v", objectType, data, err)
	}

	hashString := hashBytes.String()
	if hashString != expected {
		t.Errorf("HashObject(%q, %q) = %q, want %q", objectType, data, hashString, expected)
	}
}

func TestHashObjectEmpty(t *testing.T) {
	objectType := "blob"
	data := []byte{}
	expected := "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391"

	hashBytes, err := HashObject(objectType, data)
	if err != nil {
		t.Fatalf("HashObject(%q, %q) unexpected error: %v", objectType, data, err)
	}

	hashString := hashBytes.String()
	if hashString != expected {
		t.Errorf("HashObject(%q, %q) = %q, want %q", objectType, data, hashString, expected)
	}
}

func TestHashObjectInvalidType(t *testing.T) {
	objectType := "invalid"
	data := []byte("test")

	hashBytes, err := HashObject(objectType, data)
	if err == nil {
		t.Fatalf("HashObject(%q, %q) expected error, got nil", objectType, data)
	}

	if !hashBytes.Empty() {
		t.Errorf("HashObject(%q, %q) returned non-empty hash %q despite error", objectType, data, hashBytes.String())
	}
}
