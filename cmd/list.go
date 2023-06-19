package cmd

import (
	"fmt"
	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

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
			fmt.Print(string(content))
		}
	},
}
