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
	idFlagName = "id"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Used to get a single post",
	RunE:  getPost,
}

func getPost(cmd *cobra.Command, args []string) error {
	// Check if either id or slug flag was set
	idSet := cmd.Flags().Changed(idFlagName)
	slugSet := cmd.Flags().Changed(slugFlagName)

	if !idSet && !slugSet {
		return errors.New("either --id or --slug flag must be set")
	}

	if idSet && slugSet {
		return errors.New("cannot use both --id and --slug flags together")
	}

	if idSet {
		id, err := cmd.Flags().GetInt(idFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("Getting post with ID: %d\n", id)
	}

	if slugSet {
		slug, err := cmd.Flags().GetString(slugFlagName)
		if err != nil {
			return err
		}
		fmt.Printf("Getting post with slug: %s\n", slug)
	}

	return nil
}

func init() {
	postsCmd.AddCommand(getCmd)

	// Add flags for post retrieval
	getCmd.Flags().IntP(idFlagName, "i", 0, "ID of the post to retrieve")
	getCmd.Flags().StringP(slugFlagName, "s", "", "Slug of the post to retrieve")
}