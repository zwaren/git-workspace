package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the version number of the CLI tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
