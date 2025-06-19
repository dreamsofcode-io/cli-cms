/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// postsCmd represents the posts command
var postsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Used to manage the posts resource",
}

func init() {
	rootCmd.AddCommand(postsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}