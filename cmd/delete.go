/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	forceFlagName = "force"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove"},
	Short:   "Used to delete a post",
	RunE:    deletePost,
}

func deletePost(cmd *cobra.Command, args []string) error {
	// Check if either id or slug flag was set
	idSet := cmd.Flags().Changed(idFlagName)
	slugSet := cmd.Flags().Changed(slugFlagName)

	if !idSet && !slugSet {
		return errors.New("either --id or --slug flag must be set")
	}

	if idSet && slugSet {
		return errors.New("cannot use both --id and --slug flags together")
	}

	// Get force flag
	force, err := cmd.Flags().GetBool(forceFlagName)
	if err != nil {
		return err
	}

	// Display what will be deleted
	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("Deleting post with ID: %d\n", id)
	}

	if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("Deleting post with slug: %s\n", slug)
	}

	if force {
		fmt.Println("Force delete enabled - skipping confirmation")
	} else {
		fmt.Println("Use --force to skip confirmation")
	}

	return nil
}

func init() {
	postsCmd.AddCommand(deleteCmd)

	// Add flags for post deletion
	deleteCmd.Flags().IntP(idFlagName, "i", 0, "ID of the post to delete")
	deleteCmd.Flags().StringP(slugFlagName, "s", "", "Slug of the post to delete")
	deleteCmd.Flags().BoolP(forceFlagName, "f", false, "Force delete without confirmation")
}