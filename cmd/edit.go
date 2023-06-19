package cmd

import (
	"fmt"
	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var EditCmd = &cobra.Command{
	Use:   "edit [config|keys|hosts]",
	Short: "Edits SSH config, authorized_keys, or known_hosts file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		configPath := utils.ExpandUser(utils.SshConfigPath)

		if len(args) > 0 {
			switch args[0] {
			case "keys":
				configPath = utils.ExpandUser("~/.ssh/authorized_keys")
			case "hosts":
				configPath = utils.ExpandUser("~/.ssh/known_hosts")
			case "config":
				// Use the default configPath
			default:
				fmt.Println("Invalid argument. Use 'config' or 'authorized-keys'.")
				return
			}
		}

		command := exec.Command(editor, configPath)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}
