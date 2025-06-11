package main

import (
	"fmt"
	"os"

	"github.com/evberrypi/ssh-config/cmd"
	"github.com/evberrypi/ssh-config/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ssh-config",
	Short: "A modern SSH configuration management tool",
	Long: `ssh-config is a command-line utility that simplifies SSH configuration management.
It provides an intuitive interface for managing SSH configurations, authorized keys,
and known hosts files. Perfect for developers and system administrators who want
to streamline their SSH setup process.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version.Version,
}

func init() {
	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")

	// Add commands with aliases
	cmd.ListCmd.Aliases = []string{"ls"}
	cmd.RemoveCmd.Aliases = []string{"rm"}
	cmd.EditCmd.Aliases = []string{"e"}

	// Create help command with ? alias
	helpCmd := &cobra.Command{
		Use:     "help",
		Short:   "Help about any command",
		Aliases: []string{"?"},
		Run: func(c *cobra.Command, args []string) {
			cmd, _, err := c.Root().Find(args)
			if cmd == nil || err != nil {
				c.Printf("Unknown help topic %#q\n", args)
				c.Root().Usage()
			} else {
				cmd.InitDefaultHelpFlag()
				cmd.Help()
			}
		},
	}
	rootCmd.SetHelpCommand(helpCmd)

	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.RemoveCmd)
	rootCmd.AddCommand(cmd.EditCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
