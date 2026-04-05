package cmd

import "github.com/spf13/cobra"

var appVersion = "dev"

// SetVersion sets the application version (called from main with ldflags value).
func SetVersion(v string) {
	appVersion = v
}

var rootCmd = &cobra.Command{
	Use:   "wdxtools",
	Short: "Everyday formatting tools for the command line",
	Long:  "Lightning-fast Go ports of popular library functions that deserve their own binary.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
