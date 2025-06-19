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
	idFlagName = "id"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Used to get a single post",
	RunE:  getPost,
}

func getPost(cmd *cobra.Command, args []string) error {
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

	// Get database connection
	db, err := database.GetDatabase(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	var post *database.Post

	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("Getting post with ID: %d\n", id)
		}

		post, err = db.GetPostByID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get post: %w", err)
		}
	}

	if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("Getting post with slug: %s\n", slug)
		}

		post, err = db.GetPostBySlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("failed to get post: %w", err)
		}
	}

	// Display the post
	fmt.Printf("📝 Post Details:\n")
	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("Title: %s\n", post.Title)
	if post.Content != "" {
		fmt.Printf("Content: %s\n", post.Content)
	}
	if post.Author != "" {
		fmt.Printf("Author: %s\n", post.Author)
	}
	if post.Slug != "" {
		fmt.Printf("Slug: %s\n", post.Slug)
	}
	fmt.Printf("Created: %s\n", post.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", post.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func init() {
	postsCmd.AddCommand(getCmd)

	// Add flags for post retrieval
	getCmd.Flags().IntP(idFlagName, "i", 0, "ID of the post to retrieve")
	getCmd.Flags().StringP(slugFlagName, "s", "", "Slug of the post to retrieve")
}