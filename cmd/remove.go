package cmd

import (
	"fmt"
	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an SSH configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		configPath := utils.ExpandUser(utils.SshConfigPath)
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
