package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dreamsofcode-io/cli-cms/internal/database"
	"github.com/dreamsofcode-io/cli-cms/internal/handler/mock_handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupTestDB(t *testing.T) (*database.Database, func()) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	// Use unique filename for each test
	dbPath := filepath.Join(tempDir, fmt.Sprintf("test_%d.db", time.Now().UnixNano()))

	ctx := context.Background()
	db, err := database.New(ctx, dbPath)
	require.NoError(t, err, "Failed to create test database")

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

func TestPosts_CreatePost(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		author      string
		slug        string
		useEditor   bool
		editorResp  string
		editorErr   error
		wantErr     bool
		errContains string
	}{
		{
			name:       "Create post with editor success",
			title:      "Test Post",
			author:     "Test Author",
			slug:       "test-post",
			useEditor:  true,
			editorResp: "This is test content from editor",
			editorErr:  nil,
			wantErr:    false,
		},
		{
			name:        "Create post with editor failure",
			title:       "Test Post",
			author:      "Test Author",
			slug:        "test-post",
			useEditor:   true,
			editorResp:  "",
			editorErr:   errors.New("editor failed"),
			wantErr:     true,
			errContains: "failed to edit content",
		},
		{
			name:        "Create post with editor returning empty content",
			title:       "Test Post",
			author:      "Test Author",
			slug:        "test-post",
			useEditor:   true,
			editorResp:  "",
			editorErr:   nil,
			wantErr:     true,
			errContains: "content cannot be empty",
		},
		{
			name:      "Create post without editor",
			title:     "Test Post",
			author:    "Test Author",
			slug:      "test-post",
			useEditor: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, cleanup := setupTestDB(t)
			defer cleanup()

			mockEditor := mock_handler.NewMockTextEditor(ctrl)
			
			if tt.useEditor {
				// Set expectations for editor calls
				mockEditor.EXPECT().IsAvailable().Return(true)
				mockEditor.EXPECT().EditContentWithTemplate(tt.title, tt.author, "", false).
					Return(tt.editorResp, tt.editorErr)
			}

			// Create handler with mock editor
			handler := NewPosts(db, WithTextEditor(mockEditor))

			// Execute
			ctx := context.Background()
			createdPost, err := handler.CreatePost(ctx, tt.title, tt.author, tt.slug, tt.useEditor)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, createdPost)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, createdPost)
				assert.Equal(t, tt.title, createdPost.Title)
				assert.Equal(t, tt.author, createdPost.Author.String)
				assert.Equal(t, tt.slug, createdPost.Slug.String)
				
				if tt.useEditor {
					assert.Equal(t, tt.editorResp, createdPost.Content.String)
				}
				
				// Verify post was saved to database
				retrievedPost, err := db.GetPostByID(ctx, int(createdPost.ID))
				assert.NoError(t, err)
				assert.Equal(t, createdPost.Title, retrievedPost.Title)
			}
		})
	}
}

func TestPosts_CreatePostWithContent(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		author  string
		slug    string
		wantErr bool
	}{
		{
			name:    "Create post with content success",
			title:   "Test Post",
			content: "Test content",
			author:  "Test Author",
			slug:    "test-post",
			wantErr: false,
		},
		{
			name:    "Create post with empty content",
			title:   "Test Post",
			content: "",
			author:  "Test Author",
			slug:    "test-post",
			wantErr: false, // Empty content is allowed for this method
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, cleanup := setupTestDB(t)
			defer cleanup()

			// No need for editor mock since this method doesn't use it
			handler := NewPosts(db)

			// Execute
			ctx := context.Background()
			createdPost, err := handler.CreatePostWithContent(ctx, tt.title, tt.content, tt.author, tt.slug)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, createdPost)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, createdPost)
				assert.Equal(t, tt.title, createdPost.Title)
				assert.Equal(t, tt.content, createdPost.Content.String)
				assert.Equal(t, tt.author, createdPost.Author.String)
				assert.Equal(t, tt.slug, createdPost.Slug.String)
			}
		})
	}
}

func TestPosts_CreatePost_EditorNotAvailable(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockEditor := mock_handler.NewMockTextEditor(ctrl)
	
	// Set expectations - editor is not available
	mockEditor.EXPECT().IsAvailable().Return(false)
	mockEditor.EXPECT().GetEditorInfo().Return("nano (not found)")

	// Create handler with mock editor
	handler := NewPosts(db, WithTextEditor(mockEditor))

	// Execute
	ctx := context.Background()
	createdPost, err := handler.CreatePost(ctx, "Test Post", "Test Author", "test-post", true)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "editor not available")
	assert.Contains(t, err.Error(), "nano (not found)")
	assert.Nil(t, createdPost)
}

func TestNewPosts(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("Default configuration", func(t *testing.T) {
		handler := NewPosts(db)
		
		assert.NotNil(t, handler)
		assert.Equal(t, db, handler.db)
		assert.NotNil(t, handler.textEditor)
	})

	t.Run("With custom text editor", func(t *testing.T) {
		mockEditor := mock_handler.NewMockTextEditor(ctrl)
		
		handler := NewPosts(db, WithTextEditor(mockEditor))
		
		assert.NotNil(t, handler)
		assert.Equal(t, db, handler.db)
		assert.Equal(t, mockEditor, handler.textEditor)
	})
}

// Example of testing specific input parameters to mock
func TestPosts_CreatePost_MockInputValidation(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockEditor := mock_handler.NewMockTextEditor(ctrl)
	
	// Set strict expectations for the exact parameters
	title := "Specific Title"
	author := "Specific Author"
	slug := "specific-slug"
	expectedContent := "Editor returned this content"
	
	mockEditor.EXPECT().IsAvailable().Return(true)
	mockEditor.EXPECT().EditContentWithTemplate(
		gomock.Eq(title),    // Exact title match
		gomock.Eq(author),   // Exact author match  
		gomock.Eq(""),       // Empty existing content
		gomock.Eq(false),    // Not an update
	).Return(expectedContent, nil)

	// Create handler with mock editor
	handler := NewPosts(db, WithTextEditor(mockEditor))

	// Execute
	ctx := context.Background()
	createdPost, err := handler.CreatePost(ctx, title, author, slug, true)

	// Assert
	assert.NoError(t, err)
	require.NotNil(t, createdPost)
	assert.Equal(t, expectedContent, createdPost.Content.String)
}