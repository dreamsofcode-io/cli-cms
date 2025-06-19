/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"context"
	"fmt"
	"text/tabwriter"
	"os"

	"github.com/dreamsofcode-io/cli-cms/internal/database"
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
	ctx := context.Background()

	// Get global flags
	databaseURL, err := cmd.Flags().GetString(databaseURLFlagName)
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

	// Get database connection
	db, err := database.GetDatabase(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if verbose {
		fmt.Printf("Database URL: %s\n", databaseURL)
		fmt.Printf("Limit: %d, Offset: %d\n", limit, offset)
	}

	// Get posts from database
	posts, err := db.ListPosts(ctx, limit, offset)
	if err != nil {
		return fmt.Errorf("failed to list posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("📝 No posts found.")
		return nil
	}

	// Use tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tSLUG\tCREATED")
	fmt.Fprintln(w, "--\t-----\t------\t----\t-------")

	for _, post := range posts {
		author := post.Author
		if author == "" {
			author = "(no author)"
		}

		slug := post.Slug
		if slug == "" {
			slug = "(no slug)"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			post.ID,
			post.Title,
			author,
			slug,
			post.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	w.Flush()
	fmt.Printf("\n📊 Found %d post(s)\n", len(posts))

	return nil
}

func init() {
	postsCmd.AddCommand(listCmd)

	// Add local flags for pagination
	listCmd.Flags().IntP(limitFlagName, "l", 10, "Maximum number of posts to return")
	listCmd.Flags().IntP(offsetFlagName, "o", 0, "Number of posts to skip")
}