/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreamsofcode-io/cli-cms/internal/database"
	"github.com/dreamsofcode-io/cli-cms/internal/editor"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Short:   "Used to update a post",
	RunE:    updatePost,
}

func updatePost(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Check if either id or slug flag was set for identification
	idSet := cmd.Flags().Changed(idFlagName)
	slugSet := cmd.Flags().Changed(slugFlagName)

	if !idSet && !slugSet {
		return errors.New("either --id or --slug flag must be set to identify the post")
	}

	if idSet && slugSet {
		return errors.New("cannot use both --id and --slug flags together")
	}

	// Check if at least one update field is provided
	titleSet := cmd.Flags().Changed(titleFlagName)
	contentSet := cmd.Flags().Changed(contentFlagName)
	authorSet := cmd.Flags().Changed(authorFlagName)
	editorSet := cmd.Flags().Changed(editorFlagName)

	if !titleSet && !contentSet && !authorSet && !editorSet {
		return errors.New("at least one field must be specified to update (--title, --content, --author, or --editor)")
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

	// First, get the existing post to use for editor template
	var existingPost *database.Post
	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}
		existingPost, err = db.GetPostByID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get existing post: %w", err)
		}
	} else if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}
		existingPost, err = db.GetPostBySlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("failed to get existing post: %w", err)
		}
	}

	// Build the update struct with only changed fields
	var updates database.Post

	if titleSet {
		updates.Title, err = cmd.Flags().GetString(titleFlagName)
		if err != nil {
			return err
		}
	}

	if authorSet {
		updates.Author, err = cmd.Flags().GetString(authorFlagName)
		if err != nil {
			return err
		}
	}

	// Handle content updates
	if editorSet {
		// Use editor for content input
		useEditor, err := cmd.Flags().GetBool(editorFlagName)
		if err != nil {
			return err
		}

		if useEditor {
			if verbose {
				fmt.Printf("Opening editor for content editing...\n")
			}

			ed := editor.New()
			if !ed.IsAvailable() {
				return fmt.Errorf("editor not available: %s", ed.GetEditorInfo())
			}

			if verbose {
				fmt.Printf("Using editor: %s\n", ed.GetEditorInfo())
			}

			// Use existing post data for template
			templateTitle := existingPost.Title
			templateAuthor := existingPost.Author
			if titleSet {
				templateTitle = updates.Title
			}
			if authorSet {
				templateAuthor = updates.Author
			}

			editedContent, err := ed.EditContentWithTemplate(templateTitle, templateAuthor, existingPost.Content, true)
			if err != nil {
				return fmt.Errorf("failed to edit content: %w", err)
			}

			updates.Content = editedContent
		}
	} else if contentSet {
		// Get content from flag
		updates.Content, err = cmd.Flags().GetString(contentFlagName)
		if err != nil {
			return err
		}
	}

	var updatedPost *database.Post

	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("Updating post with ID: %d\n", id)
		}

		updatedPost, err = db.UpdatePostByID(ctx, id, updates)
		if err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}
	}

	if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("Updating post with slug: %s\n", slug)
		}

		updatedPost, err = db.UpdatePostBySlug(ctx, slug, updates)
		if err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}
	}

	// Display the updated post
	fmt.Printf("✅ Post updated successfully!\n")
	fmt.Printf("ID: %d\n", updatedPost.ID)
	fmt.Printf("Title: %s\n", updatedPost.Title)
	if updatedPost.Content != "" {
		fmt.Printf("Content: %s\n", updatedPost.Content)
	}
	if updatedPost.Author != "" {
		fmt.Printf("Author: %s\n", updatedPost.Author)
	}
	if updatedPost.Slug != "" {
		fmt.Printf("Slug: %s\n", updatedPost.Slug)
	}
	fmt.Printf("Updated: %s\n", updatedPost.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func init() {
	postsCmd.AddCommand(updateCmd)

	// Add flags for post identification
	updateCmd.Flags().IntP(idFlagName, "i", 0, "ID of the post to update")
	updateCmd.Flags().StringP(slugFlagName, "s", "", "Slug of the post to update")
	
	// Add flags for updatable fields
	updateCmd.Flags().StringP(titleFlagName, "t", "", "New title for the post")
	updateCmd.Flags().StringP(contentFlagName, "c", "", "New content for the post (ignored if --editor is used)")
	updateCmd.Flags().StringP(authorFlagName, "a", "", "New author for the post")
	updateCmd.Flags().BoolP(editorFlagName, "e", false, "Open editor for content editing")
}