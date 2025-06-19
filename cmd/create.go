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
	// Check if required title flag was set
	if !cmd.Flags().Changed(titleFlagName) {
		return errors.New("--title flag not set, must be set")
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

	// Display the post information
	fmt.Printf("Creating new post:\n")
	fmt.Printf("Title: %s\n", title)
	if content != "" {
		fmt.Printf("Content: %s\n", content)
	}
	if author != "" {
		fmt.Printf("Author: %s\n", author)
	}
	if slug != "" {
		fmt.Printf("Slug: %s\n", slug)
	}

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