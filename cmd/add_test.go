package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestAddServiceKey(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "ssh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the HOME environment variable to the temporary directory
	os.Setenv("HOME", tmpDir)

	// Mock the key that would be fetched from GitHub or GitLab
	mockKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0g+ZTxC7weoIJLUafOgrm+h..."

	// Create a test server to mock the HTTP response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, mockKey)
	}))
	defer ts.Close()

	// Replace the serviceURLs map to use the test server URL
	serviceURLs["github"] = ts.URL + "/%s.keys"

	// Mock afero filesystem
	fs := afero.NewMemMapFs()

	// Run & check ability to write comment and key
	addServiceKey("github", "testuser", fs)

	// Check if the file contains the comment and key
	authorizedKeysPath := tmpDir + "/.ssh/authorized_keys"
	aferoFs := afero.Afero{Fs: fs}
	content, err := aferoFs.ReadFile(authorizedKeysPath)
	if err != nil {
		t.Fatal(err)
	}

	comment := "# Key(s) was added from github via ssh-config"
	if !strings.Contains(string(content), comment) {
		t.Errorf("expected %q to contain %q", content, comment)
	}

	if !strings.Contains(string(content), mockKey) {
		t.Errorf("expected %q to contain %q", content, mockKey)
	}
}
