/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	databaseURLFlagName = "database-url"
	verboseFlagName     = "verbose"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cms",
	Short: "A simple application for managing blog posts",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global persistent flags available to all commands
	rootCmd.PersistentFlags().StringP(databaseURLFlagName, "d", "", "Database URL (e.g., sqlite://./blog.db)")
	rootCmd.PersistentFlags().BoolP(verboseFlagName, "v", false, "Enable verbose output")
}