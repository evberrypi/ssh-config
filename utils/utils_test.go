package utils

import (
	"os"
	"strings"
	"testing"
)

func TestExpandUser(t *testing.T) {
	home := os.Getenv("HOME")
	tests := []struct {
		input    string
		expected string
	}{
		{"~/path/to/dir", strings.Join([]string{home, "path/to/dir"}, "/")},
		{"/path/to/dir", "/path/to/dir"},
		{"~", home},
	}

	for _, test := range tests {
		result := ExpandUser(test.input)
		if result != test.expected {
			t.Errorf("ExpandUser(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
