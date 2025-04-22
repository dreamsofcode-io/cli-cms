package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "cli-cms",
	Short: "A cms tool for managing blog posts",
}

func Execute() error {
	return rootCmd.Execute()
}
