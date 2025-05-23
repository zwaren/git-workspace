package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-workspace",
	Short: "A brief description of your CLI tool",
	Long: `A longer description of your CLI tool that can span multiple lines
and likely contains examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the CLI tool!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
