/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	limitFlagName  = "limit"
	offsetFlagName = "offset"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Used to list all posts",
	RunE:  listPosts,
}

func listPosts(cmd *cobra.Command, args []string) error {
	// Get global flags
	dbURL, err := cmd.Flags().GetString(databaseURLFlagName)
	if err != nil {
		return err
	}

	verbose, err := cmd.Flags().GetBool(verboseFlagName)
	if err != nil {
		return err
	}

	// Get local flags
	limit, err := cmd.Flags().GetInt(limitFlagName)
	if err != nil {
		return err
	}

	offset, err := cmd.Flags().GetInt(offsetFlagName)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Database URL: %s\n", dbURL)
		fmt.Printf("Limit: %d, Offset: %d\n", limit, offset)
	}

	fmt.Println("Listing posts...")
	return nil
}

func init() {
	postsCmd.AddCommand(listCmd)

	// Add local flags for pagination
	listCmd.Flags().IntP(limitFlagName, "l", 10, "Maximum number of posts to return")
	listCmd.Flags().IntP(offsetFlagName, "o", 0, "Number of posts to skip")
}