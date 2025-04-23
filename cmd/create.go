/*
Copyright © 2025 Dreams of Code
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
	Short:   "Used to create a new blog post",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("posts.create called")
	},
}

func init() {
	postsCmd.AddCommand(createCmd)
}
