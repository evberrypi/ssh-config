package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/evberrypi/ssh-config/utils"
)

func TestListCmd(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "ssh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the HOME environment variable to the temporary directory
	os.Setenv("HOME", tmpDir)

	// Create .ssh directory and config file
	sshDir := filepath.Join(tmpDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}

	// Create test config content
	configContent := `Host test1
    HostName 192.168.1.1
    User user1

Host test2
    HostName 192.168.1.2
    User user2
`
	configPath := filepath.Join(sshDir, "config")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create test authorized_keys content
	keysContent := `ssh-rsa key1 user1@host1
ssh-rsa key2 user2@host2
`
	keysPath := filepath.Join(sshDir, "authorized_keys")
	if err := os.WriteFile(keysPath, []byte(keysContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Patch utils.SSHPaths to use temp dir
	oldConfig := utils.SSHPaths.Config
	oldKeys := utils.SSHPaths.AuthorizedKeys
	utils.SSHPaths.Config = configPath
	utils.SSHPaths.AuthorizedKeys = keysPath
	defer func() {
		utils.SSHPaths.Config = oldConfig
		utils.SSHPaths.AuthorizedKeys = oldKeys
	}()

	// Mock the GitHub/GitLab key response
	mockKey := "ssh-rsa mock-key user@github"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, mockKey)
	}))
	defer ts.Close()

	// Patch utils.ServiceURLs to use the test server URL
	oldURLs := utils.ServiceURLs
	utils.ServiceURLs = map[string]string{
		"github": ts.URL + "/%s.keys",
		"gitlab": ts.URL + "/%s.keys",
	}
	defer func() { utils.ServiceURLs = oldURLs }()

	tests := []struct {
		name     string
		args     []string
		expected string
		wantErr  bool
	}{
		{
			name:     "List config",
			args:     []string{"config"},
			expected: configContent,
			wantErr:  false,
		},
		{
			name:     "List keys",
			args:     []string{"keys"},
			expected: keysContent,
			wantErr:  false,
		},
		{
			name:     "List GitHub keys",
			args:     []string{"github", "testuser"},
			expected: mockKey,
			wantErr:  false,
		},
		{
			name:     "List GitLab keys",
			args:     []string{"gitlab", "testuser"},
			expected: mockKey,
			wantErr:  false,
		},
		{
			name:     "Invalid platform",
			args:     []string{"invalid", "testuser"},
			expected: "Unknown platform. Use 'github' or 'gitlab'.",
			wantErr:  false,
		},
		{
			name:     "Invalid argument",
			args:     []string{"invalid"},
			expected: "Invalid argument. Use 'config', 'keys', 'github [username]' or 'gitlab [username]'.",
			wantErr:  false,
		},
		{
			name:     "No arguments",
			args:     []string{},
			expected: "Invalid number of arguments. Use 'ssh-config list --help' for usage information.",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := ListCmd
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			// Execute the command
			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check the output
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("ListCmd output = %v, want %v", output, tt.expected)
			}
		})
	}
}

func TestListCmdWithMissingFiles(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "ssh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the HOME environment variable to the temporary directory
	os.Setenv("HOME", tmpDir)

	// Create .ssh directory without any files
	sshDir := filepath.Join(tmpDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}

	// Patch utils.SSHPaths to use temp dir
	oldConfig := utils.SSHPaths.Config
	oldKeys := utils.SSHPaths.AuthorizedKeys
	utils.SSHPaths.Config = filepath.Join(sshDir, "config")
	utils.SSHPaths.AuthorizedKeys = filepath.Join(sshDir, "authorized_keys")
	defer func() {
		utils.SSHPaths.Config = oldConfig
		utils.SSHPaths.AuthorizedKeys = oldKeys
	}()

	tests := []struct {
		name     string
		args     []string
		expected string
		wantErr  bool
	}{
		{
			name:     "List missing config",
			args:     []string{"config"},
			expected: "Error:",
			wantErr:  false,
		},
		{
			name:     "List missing keys",
			args:     []string{"keys"},
			expected: "Error:",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := ListCmd
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			// Execute the command
			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check the output
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("ListCmd output = %v, want %v", output, tt.expected)
			}
		})
	}
}
