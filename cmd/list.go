// Package cmd provides command-line interfaces for interacting with SSH configurations.
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
)

// ListCmd represents the Cobra command for listing SSH configuration of ~/.ssh/config
// or the public keys on gitlab.com and github.com for a specific user.
var ListCmd = &cobra.Command{
	Use:   "list [config|keys|github|gitlab] [username]",
	Short: "List SSH configurations or fetch SSH keys from GitHub/GitLab",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			platform := args[0]
			username := args[1]
			url := ""

			switch platform {
			case "github":
				url = "https://github.com/" + username + ".keys"
			case "gitlab":
				url = "https://gitlab.com/" + username + ".keys"
			default:
				fmt.Println("Unknown platform. Use 'github' or 'gitlab'.")
				return
			}

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error fetching keys:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				fmt.Println(string(body))
			} else {
				fmt.Println("Error fetching keys. HTTP status:", resp.Status)
			}
		} else if len(args) == 1 {
			switch args[0] {
			case "keys":
				keysPath := utils.ExpandUser("~/.ssh/authorized_keys")
				content, err := os.ReadFile(keysPath)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println(string(content))
			case "config":
				configPath := utils.ExpandUser(utils.SSHPaths.Config)
				content, err := os.ReadFile(configPath)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(content)) // Write to the command's output stream
			default:
				fmt.Println("Invalid argument. Use 'config', 'keys', 'github [username]' or 'gitlab [username]'.")
			}
		} else {
			fmt.Println("Invalid number of arguments. Use 'ssh-config list --help' for usage information.")
		}
	},
}
