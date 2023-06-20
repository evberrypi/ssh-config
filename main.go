package main

import (
	"github.com/evberrypi/ssh-config/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ssh-config",
	Short: "Manage SSH configurations.",
}

// Basic main function. Please see 'cmd' and 'utils' directories more of what is going on
func main() {
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.RemoveCmd)
	rootCmd.AddCommand(cmd.EditCmd)
	rootCmd.Execute()
}
