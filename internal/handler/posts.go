//go:generate mockgen -source=posts.go -destination=mock_handler/posts.go

package handler

import (
	"context"
	"fmt"

	"github.com/dreamsofcode-io/cli-cms/internal/database"
	"github.com/dreamsofcode-io/cli-cms/internal/editor"
)

// TextEditor interface defines the methods needed for text editing
type TextEditor interface {
	EditContentWithTemplate(title, author, existingContent string, isUpdate bool) (string, error)
	IsAvailable() bool
	GetEditorInfo() string
}

// Posts handles post-related operations
type Posts struct {
	db         *database.Database
	textEditor TextEditor
}

// Option defines a function type for configuring Posts
type Option func(*Posts)

// WithTextEditor returns an Option to configure the text editor
func WithTextEditor(editor TextEditor) Option {
	return func(p *Posts) {
		p.textEditor = editor
	}
}

// NewPosts creates a new Posts handler with optional configuration
func NewPosts(db *database.Database, opts ...Option) *Posts {
	res := &Posts{
		db:         db,
		textEditor: editor.New(),
	}
	
	for _, opt := range opts {
		opt(res)
	}
	
	return res
}

// CreatePost creates a new blog post
func (p *Posts) CreatePost(ctx context.Context, title, author, slug string, useEditor bool) (*database.Post, error) {
	var content string
	
	if useEditor {
		if !p.textEditor.IsAvailable() {
			return nil, fmt.Errorf("editor not available: %s", p.textEditor.GetEditorInfo())
		}
		
		editedContent, err := p.textEditor.EditContentWithTemplate(title, author, "", false)
		if err != nil {
			return nil, fmt.Errorf("failed to edit content: %w", err)
		}
		
		content = editedContent
		
		if content == "" {
			return nil, fmt.Errorf("content cannot be empty when using editor")
		}
	}
	
	// Create the post using helper function
	post := database.CreatePostFromInput(title, content, author, slug)
	
	createdPost, err := p.db.CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	
	return createdPost, nil
}

// CreatePostWithContent creates a new blog post with provided content
func (p *Posts) CreatePostWithContent(ctx context.Context, title, content, author, slug string) (*database.Post, error) {
	// Create the post using helper function
	post := database.CreatePostFromInput(title, content, author, slug)
	
	createdPost, err := p.db.CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	
	return createdPost, nil
}