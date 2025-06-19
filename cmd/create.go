/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Used to create a new post",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("post create called")
		return nil
	},
}

func init() {
	postsCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}