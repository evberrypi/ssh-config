package utils

import (
	"os"
	"strings"
)

var SshConfigPath = "~/.ssh/config"

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
