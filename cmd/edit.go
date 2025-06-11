// Package cmd provides command-line interfaces for interacting with SSH configurations.
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
)

var execCommand = exec.Command
var getenvFunc = os.Getenv

var EditCmd = &cobra.Command{
	Use:   "edit [config|keys|hosts]",
	Short: "Edits SSH config, authorized_keys, or known_hosts file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		editor := getenvFunc("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		configPath := utils.ExpandUser(utils.SSHPaths.Config)

		if len(args) > 0 {
			switch args[0] {
			case "keys":
				configPath = utils.ExpandUser("~/.ssh/authorized_keys")
			case "hosts":
				configPath = utils.ExpandUser("~/.ssh/known_hosts")
			case "config":
				// Use the default configPath
			default:
				cmd.Println("Invalid argument. Use 'config' or 'keys'.")
				return
			}
		}

		command := execCommand(editor, configPath)
		command.Stdin = os.Stdin
		command.Stdout = cmd.OutOrStdout() // This line captures the output
		command.Stderr = os.Stderr
		err := command.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}
