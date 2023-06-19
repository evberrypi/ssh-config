// Package cmd provides command-line interfaces for interacting with SSH configurations.
package cmd

import (
	"fmt"
	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

// ListCmd represents the Cobra command for listing SSH configuration of ~/.ssh/config or the public keys on gitlab.com and github.com for a specific user.
var ListCmd = &cobra.Command{
	Use:   "list [github-keys|gitlab-keys] [username]",
	Short: "List SSH configurations or fetch SSH keys from GitHub/GitLab",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			platform := args[0]
			username := args[1]
			url := ""

			switch platform {
			case "github-keys":
				url = "https://github.com/" + username + ".keys"
			case "gitlab-keys":
				url = "https://gitlab.com/" + username + ".keys"
			default:
				fmt.Println("Unknown platform. Use 'github-keys' or 'gitlab-keys'.")
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
		} else {
			configPath := utils.ExpandUser(utils.SshConfigPath)
			content, err := os.ReadFile(configPath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			cmd.Print(string(content))
		}
	},
}
