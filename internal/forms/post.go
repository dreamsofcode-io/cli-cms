package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dreamsofcode-io/cli-cms/internal/database"
	"github.com/dreamsofcode-io/cli-cms/internal/editor"
)

// PostFormData holds the data collected from the post creation form
type PostFormData struct {
	Title   string
	Content string
	Author  string
	Slug    string
	Confirm bool
}

// NewPostForm creates an interactive form for creating a blog post
// If useEditor is true, it will open an editor for content input after the form
func NewPostForm(initialData PostFormData, useEditor bool) (*PostFormData, error) {
	var data PostFormData
	
	// Pre-populate with initial data if provided
	data.Title = initialData.Title
	data.Content = initialData.Content  
	data.Author = initialData.Author
	data.Slug = initialData.Slug

	// Create form fields
	titleInput := huh.NewInput().
		Value(&data.Title).
		Title("Post Title").
		Description("The title of your blog post").
		Placeholder("Enter a catchy title").
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("title is required")
			}
			return nil
		})

	// Create content input based on useEditor flag
	var contentInput huh.Field
	if useEditor {
		contentInput = huh.NewInput().
			Value(&data.Content).
			Title("Content (Editor)").
			Description("Content will be edited in your text editor after this form").
			Placeholder("[Will open editor for content]")
	} else {
		// Fallback to text area if editor not available
		contentInput = huh.NewText().
			Value(&data.Content).
			Title("Post Content").
			Description("Enter your post content").
			Placeholder("Write your post content here...")
	}

	authorInput := huh.NewInput().
		Value(&data.Author).
		Title("Author").
		Description("The author of this blog post").
		Placeholder("Enter author name")

	slugInput := huh.NewInput().
		Value(&data.Slug).
		Title("URL Slug").
		Description("The URL-friendly version of the title (optional)").
		Placeholder("my-awesome-post")

	confirmInput := huh.NewConfirm().
		Title("Create this blog post?").
		Description("Review your post details above and confirm creation").
		Value(&data.Confirm)

	// Create form with all fields
	form := huh.NewForm(
		huh.NewGroup(
			titleInput,
			contentInput,
			authorInput,
			slugInput,
		),
		huh.NewGroup(
			confirmInput,
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		return nil, err
	}

	// If user didn't confirm, return early
	if !data.Confirm {
		return nil, ErrUserCancelled
	}

	// Open editor for content if enabled
	if useEditor {
		ed := editor.New()
		if !ed.IsAvailable() {
			return nil, fmt.Errorf("editor not available: %s", ed.GetEditorInfo())
		}

		editedContent, err := ed.EditContentWithTemplate(data.Title, data.Author, data.Content, false)
		if err != nil {
			return nil, fmt.Errorf("failed to edit content: %w", err)
		}

		data.Content = editedContent

		if data.Content == "" {
			return nil, fmt.Errorf("content cannot be empty")
		}
	}

	// Auto-generate slug if not provided
	if strings.TrimSpace(data.Slug) == "" {
		data.Slug = generateSlug(data.Title)
	}

	return &data, nil
}

// ToPost converts form data to a database Post
func (data *PostFormData) ToPost() database.Post {
	return database.CreatePostFromInput(
		data.Title,
		data.Content,
		data.Author,
		data.Slug,
	)
}

// generateSlug creates a URL-friendly slug from a title
func generateSlug(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove special characters (keep only letters, numbers, and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	// Clean up multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}