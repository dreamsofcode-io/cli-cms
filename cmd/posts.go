package cmd

import (
	"github.com/spf13/cobra"
)

var postsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Used to manage the blog posts resource",
}

func init() {
	rootCmd.AddCommand(postsCmd)
}
