package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandUser(t *testing.T) {
	home := os.Getenv("HOME")
	tests := []struct {
		input    string
		expected string
	}{
		{"~/path/to/dir", filepath.Join(home, "path/to/dir")},
		{"/path/to/dir", "/path/to/dir"},
		{"~", home},
		{"~/.ssh/config", filepath.Join(home, ".ssh/config")},
		{"no-tilde", "no-tilde"},
	}

	for _, test := range tests {
		result := ExpandUser(test.input)
		if result != test.expected {
			t.Errorf("ExpandUser(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestEnsureSSHDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ssh-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Test creating .ssh directory
	err = EnsureSSHDirectory()
	if err != nil {
		t.Errorf("EnsureSSHDirectory() error = %v", err)
	}

	// Verify directory exists with correct permissions
	sshDir := filepath.Join(tmpDir, ".ssh")
	info, err := os.Stat(sshDir)
	if err != nil {
		t.Errorf("Failed to stat .ssh directory: %v", err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("Directory permissions = %v; want 0700", info.Mode().Perm())
	}
}

func TestEnsureFileExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ssh-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name    string
		path    string
		perm    os.FileMode
		wantErr bool
	}{
		{
			name:    "Create new file",
			path:    "~/.ssh/config",
			perm:    0644,
			wantErr: false,
		},
		{
			name:    "Create new authorized_keys",
			path:    "~/.ssh/authorized_keys",
			perm:    0600,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureFileExists(tt.path, tt.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureFileExists() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify file exists with correct permissions
			expandedPath := ExpandUser(tt.path)
			info, err := os.Stat(expandedPath)
			if err != nil {
				t.Errorf("Failed to stat file: %v", err)
			}
			if info.Mode().Perm() != tt.perm {
				t.Errorf("File permissions = %v; want %v", info.Mode().Perm(), tt.perm)
			}
		})
	}
}

func TestGetDefaultEditor(t *testing.T) {
	// Test with EDITOR set
	oldEditor := os.Getenv("EDITOR")
	os.Setenv("EDITOR", "nano")
	defer os.Setenv("EDITOR", oldEditor)

	if got := GetDefaultEditor(); got != "nano" {
		t.Errorf("GetDefaultEditor() = %v; want nano", got)
	}

	// Test with EDITOR unset
	os.Unsetenv("EDITOR")
	if got := GetDefaultEditor(); got != "vim" {
		t.Errorf("GetDefaultEditor() = %v; want vim", got)
	}
}

func TestFormatSSHConfig(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		hostname     string
		user         string
		identityFile string
		extraArgs    map[string]string
		want         string
	}{
		{
			name:         "Basic config",
			host:         "test",
			hostname:     "example.com",
			user:         "user",
			identityFile: "~/.ssh/id_rsa",
			extraArgs:    nil,
			want:         "Host test\n    HostName example.com\n    User user\n    IdentityFile ~/.ssh/id_rsa\n",
		},
		{
			name:         "Config with extra args",
			host:         "test",
			hostname:     "example.com",
			user:         "user",
			identityFile: "~/.ssh/id_rsa",
			extraArgs: map[string]string{
				"Port":         "2222",
				"ForwardAgent": "yes",
			},
			want: "Host test\n    HostName example.com\n    User user\n    IdentityFile ~/.ssh/id_rsa\n    Port 2222\n    ForwardAgent yes\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSSHConfig(tt.host, tt.hostname, tt.user, tt.identityFile, tt.extraArgs)
			if got != tt.want {
				t.Errorf("FormatSSHConfig() = %v; want %v", got, tt.want)
			}
		})
	}
}
