package main

import (
	"github.com/spf13/cobra"
	"github.com/evberrypi/ssh/config/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "ssh-config",
	Short: "Manage SSH configurations.",
}

func main() {
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.RemoveCmd)
	rootCmd.AddCommand(cmd.EditCmd)
	rootCmd.Execute()
}
