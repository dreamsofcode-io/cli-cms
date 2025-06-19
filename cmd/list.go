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
	"github.com/dreamsofcode-io/cli-cms/internal/ui"
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

	// Display header
	ui.Header("Posts")
	
	// Use tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, ui.HighlightString("ID\tTITLE\tAUTHOR\tSLUG\tCREATED"))
	fmt.Fprintln(w, ui.SubtleString("--\t-----\t------\t----\t-------"))

	for _, post := range posts {
		author := "(no author)"
		if post.Author.Valid {
			author = post.Author.String
		}

		slug := "(no slug)"
		if post.Slug.Valid {
			slug = post.Slug.String
		}

		created := "(no date)"
		if post.CreatedAt.Valid {
			created = post.CreatedAt.Time.Format("2006-01-02 15:04")
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			post.ID,
			ui.HighlightString(post.Title),
			author,
			ui.LinkString(slug),
			ui.SubtleString(created),
		)
	}

	w.Flush()
	fmt.Printf("\n")
	ui.PrintInfo("Found %d post(s)\n", len(posts))

	return nil
}

func init() {
	postsCmd.AddCommand(listCmd)

	// Add local flags for pagination
	listCmd.Flags().IntP(limitFlagName, "l", 10, "Maximum number of posts to return")
	listCmd.Flags().IntP(offsetFlagName, "o", 0, "Number of posts to skip")
}