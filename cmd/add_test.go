package cmd

import (
	"os"
	"strings"
	"testing"
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
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0g+ZTxC7weoIJLUafOgrm+h..."

	// Run the function - for this test example, we'll just check if it's able to write the comment and key correctly
	addServiceKey("github", "testuser")

	// Check if the file contains the comment and key
	authorizedKeysPath := tmpDir + "/.ssh/authorized_keys"
	content, err := os.ReadFile(authorizedKeysPath)
	if err != nil {
		t.Fatal(err)
	}

	comment := "# This key was added via github via the ssh-config tool"
	if !strings.Contains(string(content), comment) {
		t.Errorf("expected %q to contain %q", content, comment)
	}

	if !strings.Contains(string(content), key) {
		t.Errorf("expected %q to contain %q", content, key)
	}
}
