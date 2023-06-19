package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
)

func TestRemoveCmd(t *testing.T) {
	// Create a temporary file as mock ssh config
	tmpfile, err := os.CreateTemp("", "example.*.config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write sample data to the file
	configContent := `Host test1
    HostName 192.168.1.1
    User user1

Host test2
    HostName 192.168.1.2
    User user2
`
	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Mock SshConfigPath in utils package
	utils.SshConfigPath = tmpfile.Name()

	// Execute RemoveCmd
	cmd := &cobra.Command{}
	cmd.AddCommand(RemoveCmd)
	cmd.SetArgs([]string{"remove", "test1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Check the content of the file
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := `Host test2
    HostName 192.168.1.2
    User user2
`

	if strings.TrimSpace(string(content)) != strings.TrimSpace(expectedContent) {
		t.Errorf("Expected content %q, but got %q", expectedContent, content)
	}
}
