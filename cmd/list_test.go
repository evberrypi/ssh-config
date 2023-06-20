package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
)

func TestListCmd(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write sample content to temporary file
	content := []byte("Host test1\n    HostName 192.168.1.1\n    User user1\n\nHost test2\n    HostName 192.168.1.2\n    User user2\n")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Mock SshConfigPath in utils package
	utils.SshConfigPath = tmpfile.Name()

	// Buffer to capture stdout
	var buf bytes.Buffer

	// Run the list command
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List SSH configurations or fetch SSH keys from GitHub/GitLab",
		Run:   ListCmd.Run,
	}
	// Here, pass "config" argument to command.
	cmd.SetArgs([]string{"config"})
	cmd.SetOut(&buf)
	cmd.Execute()

	// Check the output
	expected := "Host test1\n    HostName 192.168.1.1\n    User user1\n\nHost test2\n    HostName 192.168.1.2\n    User user2\n"
	actual := buf.String()
	if actual != expected {
		t.Errorf("Expected output %q, but got %q", expected, actual)
	}
}
