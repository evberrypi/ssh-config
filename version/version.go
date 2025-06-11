package version

import (
	"fmt"
)

var (
	// Version is the current version of ssh-config
	Version = "1.0.0"
	// BuildTime is the time the binary was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// String returns the full version string
func String() string {
	return fmt.Sprintf("ssh-config v%s (build: %s, commit: %s)", Version, BuildTime, GitCommit)
}
