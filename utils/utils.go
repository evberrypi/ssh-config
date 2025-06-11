package utils

import (
	"fmt"
	"os"
	"strings"
)

const (
	// DefaultSSHConfigPath is the default path to the SSH config file
	DefaultSSHConfigPath = "~/.ssh/config"
	// DefaultAuthorizedKeysPath is the default path to the authorized_keys file
	DefaultAuthorizedKeysPath = "~/.ssh/authorized_keys"
	// DefaultKnownHostsPath is the default path to the known_hosts file
	DefaultKnownHostsPath = "~/.ssh/known_hosts"
)

// SSHPaths contains all the relevant SSH file paths
var SSHPaths = struct {
	Config         string
	AuthorizedKeys string
	KnownHosts     string
}{
	Config:         DefaultSSHConfigPath,
	AuthorizedKeys: DefaultAuthorizedKeysPath,
	KnownHosts:     DefaultKnownHostsPath,
}

// ExpandUser expands the tilde (~) in a file path to the user's home directory.
// If the path is "~", it returns the value of the HOME environment variable.
// If the path starts with "~/", it replaces the tilde with the home directory path.
// Otherwise, it returns the path as-is.
func ExpandUser(path string) string {
	if path == "~" {
		return os.Getenv("HOME")
	} else if len(path) >= 2 && path[:2] == "~/" {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	return path
}

// EnsureSSHDirectory ensures that the .ssh directory exists and has the correct permissions
func EnsureSSHDirectory() error {
	sshDir := ExpandUser("~/.ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}
	return nil
}

// EnsureFileExists ensures that a file exists and has the correct permissions
func EnsureFileExists(path string, perm os.FileMode) error {
	expandedPath := ExpandUser(path)
	if err := EnsureSSHDirectory(); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		// Create file with correct permissions
		file, err := os.OpenFile(expandedPath, os.O_CREATE|os.O_WRONLY, perm)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		file.Close()
	} else if err != nil {
		return fmt.Errorf("failed to check file %s: %w", path, err)
	}

	// Ensure correct permissions
	if err := os.Chmod(expandedPath, perm); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", path, err)
	}

	return nil
}

// GetDefaultEditor returns the default editor to use
func GetDefaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "vim" // Default fallback
}

// FormatSSHConfig formats an SSH config block with proper indentation
func FormatSSHConfig(host, hostname, user, identityFile string, extraArgs map[string]string) string {
	config := fmt.Sprintf("Host %s\n    HostName %s\n    User %s\n    IdentityFile %s\n",
		host, hostname, user, identityFile)

	for key, value := range extraArgs {
		config += fmt.Sprintf("    %s %s\n", key, value)
	}

	return config
}
