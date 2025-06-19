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
	titleFlagName   = "title"
	contentFlagName = "content"
	authorFlagName  = "author"
	slugFlagName    = "slug"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Used to create a new post",
	RunE:    createPost,
}

func createPost(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Check if required title flag was set
	if !cmd.Flags().Changed(titleFlagName) {
		return errors.New("--title flag not set, must be set")
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

	// Get flag values
	title, err := cmd.Flags().GetString(titleFlagName)
	if err != nil {
		return err
	}

	content, err := cmd.Flags().GetString(contentFlagName)
	if err != nil {
		return err
	}

	author, err := cmd.Flags().GetString(authorFlagName)
	if err != nil {
		return err
	}

	slug, err := cmd.Flags().GetString(slugFlagName)
	if err != nil {
		return err
	}

	// Get database connection
	db, err := database.GetDatabase(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if verbose {
		fmt.Printf("Database URL: %s\n", databaseURL)
		fmt.Printf("Creating new post...\n")
	}

	// Create the post
	post := database.Post{
		Title:   title,
		Content: content,
		Author:  author,
		Slug:    slug,
	}

	createdPost, err := db.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	// Display the created post information
	fmt.Printf("✅ Post created successfully!\n")
	fmt.Printf("ID: %d\n", createdPost.ID)
	fmt.Printf("Title: %s\n", createdPost.Title)
	if createdPost.Content != "" {
		fmt.Printf("Content: %s\n", createdPost.Content)
	}
	if createdPost.Author != "" {
		fmt.Printf("Author: %s\n", createdPost.Author)
	}
	if createdPost.Slug != "" {
		fmt.Printf("Slug: %s\n", createdPost.Slug)
	}
	fmt.Printf("Created: %s\n", createdPost.CreatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func init() {
	postsCmd.AddCommand(createCmd)

	// Add flags for post creation
	createCmd.Flags().StringP(titleFlagName, "t", "", "Title of the post (required)")
	createCmd.Flags().StringP(contentFlagName, "c", "", "Content of the post")
	createCmd.Flags().StringP(authorFlagName, "a", "", "Author of the post")
	createCmd.Flags().StringP(slugFlagName, "s", "", "URL slug for the post")
}