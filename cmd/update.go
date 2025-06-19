/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"errors"
	"fmt"

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

	if !titleSet && !contentSet && !authorSet {
		return errors.New("at least one field must be specified to update (--title, --content, or --author)")
	}

	// Display which post is being updated
	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("Updating post with ID: %d\n", id)
	}

	if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("Updating post with slug: %s\n", slug)
	}

	// Show what fields are being updated
	fmt.Println("Fields to update:")

	if titleSet {
		title, err := cmd.Flags().GetString(titleFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("  Title: %s\n", title)
	}

	if contentSet {
		content, err := cmd.Flags().GetString(contentFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("  Content: %s\n", content)
	}

	if authorSet {
		author, err := cmd.Flags().GetString(authorFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("  Author: %s\n", author)
	}

	return nil
}

func init() {
	postsCmd.AddCommand(updateCmd)

	// Add flags for post identification
	updateCmd.Flags().IntP(idFlagName, "i", 0, "ID of the post to update")
	updateCmd.Flags().StringP(slugFlagName, "s", "", "Slug of the post to update")
	
	// Add flags for updatable fields
	updateCmd.Flags().StringP(titleFlagName, "t", "", "New title for the post")
	updateCmd.Flags().StringP(contentFlagName, "c", "", "New content for the post")
	updateCmd.Flags().StringP(authorFlagName, "a", "", "New author for the post")
}