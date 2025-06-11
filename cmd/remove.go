// Package cmd provides command-line interfaces for interacting with SSH configurations.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
)

// RemoveCmd represents the Cobra command for removing a host from the SSH configuration file ~/.ssh/config.
var RemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a SSH host from the SSH configuration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		configPath := utils.ExpandUser(utils.SSHPaths.Config)
		content, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		lines := strings.Split(string(content), "\n")
		newContent := ""
		skip := false
		for _, line := range lines {
			if strings.HasPrefix(line, "Host "+name) {
				skip = true
			} else if skip && strings.HasPrefix(line, "    ") {
				continue
			} else {
				skip = false
				newContent += line + "\n"
			}
		}

		err = os.WriteFile(configPath, []byte(newContent), 0644)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}
