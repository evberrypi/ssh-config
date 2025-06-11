package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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

	// Patch utils.ServiceURLs to use the test server URL
	oldURLs := utils.ServiceURLs
	utils.ServiceURLs = map[string]string{
		"github": ts.URL + "/%s.keys",
		"gitlab": ts.URL + "/%s.keys",
	}
	defer func() { utils.ServiceURLs = oldURLs }()

	// Mock afero filesystem
	fs := afero.NewMemMapFs()

	// Create .ssh directory and authorized_keys file
	sshDir := filepath.Join(tmpDir, ".ssh")
	if err := fs.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}

	authorizedKeysPath := filepath.Join(sshDir, "authorized_keys")
	if _, err := fs.Create(authorizedKeysPath); err != nil {
		t.Fatal(err)
	}

	// Run & check ability to write comment and key
	addServiceKey("github", "testuser", fs)

	// Check if the file contains the comment and key
	aferoFs := afero.Afero{Fs: fs}
	content, err := aferoFs.ReadFile(authorizedKeysPath)
	if err != nil {
		t.Fatal(err)
	}

	comment := "# Keys added from github user testuser via ssh-config"
	if !strings.Contains(string(content), comment) {
		t.Errorf("expected %q to contain %q", content, comment)
	}

	if !strings.Contains(string(content), mockKey) {
		t.Errorf("expected %q to contain %q", content, mockKey)
	}
}

func TestAddServiceKeyErrors(t *testing.T) {
	tests := []struct {
		name     string
		service  string
		username string
		fs       afero.Fs
		wantErr  bool
	}{
		{
			name:     "Invalid service",
			service:  "invalid",
			username: "testuser",
			fs:       afero.NewMemMapFs(),
			wantErr:  true,
		},
		{
			name:     "Empty username",
			service:  "github",
			username: "",
			fs:       afero.NewMemMapFs(),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addServiceKey(tt.service, tt.username, tt.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("addServiceKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddConfig(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "ssh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the HOME environment variable to the temporary directory
	os.Setenv("HOME", tmpDir)

	// Patch utils.SSHPaths to use temp dir
	sshDir := filepath.Join(tmpDir, ".ssh")
	oldConfig := utils.SSHPaths.Config
	oldKeys := utils.SSHPaths.AuthorizedKeys
	utils.SSHPaths.Config = filepath.Join(sshDir, "config")
	utils.SSHPaths.AuthorizedKeys = filepath.Join(sshDir, "authorized_keys")
	defer func() {
		utils.SSHPaths.Config = oldConfig
		utils.SSHPaths.AuthorizedKeys = oldKeys
	}()

	// Patch promptForExtraArgs to avoid stdin
	oldPrompt := promptForExtraArgs
	promptForExtraArgs = func(reader *bufio.Reader) (map[string]string, error) {
		return map[string]string{}, nil
	}
	defer func() { promptForExtraArgs = oldPrompt }()

	// Set up test config
	configOptions = ConfigOptions{
		HostName:  "test",
		IPAddress: "192.168.1.1",
		Username:  "testuser",
		SSHKey:    "~/.ssh/id_rsa",
	}

	// Create a buffer to capture stdout
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	// Run the config command
	err = runConfigCmd(cmd, []string{})
	if err != nil {
		t.Fatalf("runConfigCmd() error = %v", err)
	}

	// Read the config file
	configPath := filepath.Join(sshDir, "config")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Check if the config contains the expected content
	expected := fmt.Sprintf("Host %s\n    HostName %s\n    User %s\n    IdentityFile %s\n",
		configOptions.HostName,
		configOptions.IPAddress,
		configOptions.Username,
		configOptions.SSHKey,
	)

	if !strings.Contains(string(content), expected) {
		t.Errorf("Config content = %v; want %v", string(content), expected)
	}
}

func TestAddConfigWithExtraArgs(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "ssh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the HOME environment variable to the temporary directory
	os.Setenv("HOME", tmpDir)

	// Patch utils.SSHPaths to use temp dir
	sshDir := filepath.Join(tmpDir, ".ssh")
	oldConfig := utils.SSHPaths.Config
	oldKeys := utils.SSHPaths.AuthorizedKeys
	utils.SSHPaths.Config = filepath.Join(sshDir, "config")
	utils.SSHPaths.AuthorizedKeys = filepath.Join(sshDir, "authorized_keys")
	defer func() {
		utils.SSHPaths.Config = oldConfig
		utils.SSHPaths.AuthorizedKeys = oldKeys
	}()

	// Patch promptForExtraArgs to avoid stdin
	oldPrompt := promptForExtraArgs
	promptForExtraArgs = func(reader *bufio.Reader) (map[string]string, error) {
		return map[string]string{
			"Port":         "2222",
			"ForwardAgent": "yes",
		}, nil
	}
	defer func() { promptForExtraArgs = oldPrompt }()

	// Set up test config with extra args
	configOptions = ConfigOptions{
		HostName:  "test",
		IPAddress: "192.168.1.1",
		Username:  "testuser",
		SSHKey:    "~/.ssh/id_rsa",
	}

	// Create a buffer to capture stdout
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	// Run the config command
	err = runConfigCmd(cmd, []string{})
	if err != nil {
		t.Fatalf("runConfigCmd() error = %v", err)
	}

	// Read the config file
	configPath := filepath.Join(sshDir, "config")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Check if the config contains all expected content
	expected := fmt.Sprintf("Host %s\n    HostName %s\n    User %s\n    IdentityFile %s\n    Port 2222\n    ForwardAgent yes\n",
		configOptions.HostName,
		configOptions.IPAddress,
		configOptions.Username,
		configOptions.SSHKey,
	)

	if !strings.Contains(string(content), expected) {
		t.Errorf("Config content = %v; want %v", string(content), expected)
	}
}
