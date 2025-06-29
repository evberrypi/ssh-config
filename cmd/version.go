package cmd

import (
	"github.com/evberrypi/ssh-config/version"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the version number of ssh-config`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version.String())
	},
}
