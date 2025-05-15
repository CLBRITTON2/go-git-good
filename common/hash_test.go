package common

import (
	"testing"
)

func TestHashObject(t *testing.T) {
	tests := []struct {
		name       string
		objectType string
		data       []byte
		expected   string
	}{
		{
			// Data from echo -n "test" | git hash-object --stdin
			name:       "blob with simple text",
			objectType: "blob",
			data:       []byte("test"),
			expected:   "30d74d258442c7c65512eafab474568dd706c430",
		},
		{
			name:       "blob with empty data",
			objectType: "blob",
			data:       []byte{},
			expected:   "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashBytes := HashObject(test.objectType, test.data)
			hashString := hashBytes.String()
			if hashString != test.expected {
				t.Errorf("HashObject(%q, %q) = %q, want %q", test.objectType, test.data, hashString, test.expected)
			}
		})
	}
}
