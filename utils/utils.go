package utils

import (
	"os"
	"strings"
)

const SshConfigPath = "~/.ssh/config"

func ExpandUser(path string) string {
	if path[:2] == "~/" {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	return path
}
