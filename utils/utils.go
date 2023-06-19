package utils

import (
	"os"
	"strings"
)

var SshConfigPath = "~/.ssh/config"

func ExpandUser(path string) string {
	if path == "~" {
		return os.Getenv("HOME")
	} else if len(path) >= 2 && path[:2] == "~/" {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	return path
}
