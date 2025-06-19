/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Short:   "Used to update a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("post update called")
		return nil
	},
}

func init() {
	postsCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}