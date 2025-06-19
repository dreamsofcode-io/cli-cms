/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove"},
	Short:   "Used to delete a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("post delete called")
		return nil
	},
}

func init() {
	postsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}