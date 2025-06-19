/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreamsofcode-io/cli-cms/internal/database"
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
	ctx := context.Background()

	// Check if either id or slug flag was set
	idSet := cmd.Flags().Changed(idFlagName)
	slugSet := cmd.Flags().Changed(slugFlagName)

	if !idSet && !slugSet {
		return errors.New("either --id or --slug flag must be set")
	}

	if idSet && slugSet {
		return errors.New("cannot use both --id and --slug flags together")
	}

	// Get database URL from global flag
	databaseURL, err := cmd.Flags().GetString(databaseURLFlagName)
	if err != nil {
		return err
	}

	// Get verbose flag
	verbose, err := cmd.Flags().GetBool(verboseFlagName)
	if err != nil {
		return err
	}

	// Get force flag
	force, err := cmd.Flags().GetBool(forceFlagName)
	if err != nil {
		return err
	}

	// Get database connection
	db, err := database.GetDatabase(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("Deleting post with ID: %d\n", id)
		}

		if !force {
			fmt.Printf("⚠️  Are you sure you want to delete post ID %d? Use --force to skip this confirmation.\n", id)
			return nil
		}

		err = db.DeletePostByID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete post: %w", err)
		}

		fmt.Printf("✅ Post with ID %d deleted successfully!\n", id)
	}

	if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("Deleting post with slug: %s\n", slug)
		}

		if !force {
			fmt.Printf("⚠️  Are you sure you want to delete post with slug '%s'? Use --force to skip this confirmation.\n", slug)
			return nil
		}

		err = db.DeletePostBySlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("failed to delete post: %w", err)
		}

		fmt.Printf("✅ Post with slug '%s' deleted successfully!\n", slug)
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