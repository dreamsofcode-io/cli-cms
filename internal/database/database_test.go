package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*Database, func()) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	// Use unique filename for each test
	dbPath := filepath.Join(tempDir, fmt.Sprintf("test_%d.db", time.Now().UnixNano()))

	ctx := context.Background()
	db, err := New(ctx, dbPath)
	require.NoError(t, err, "Failed to create test database")

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
		wantErr     bool
	}{
		{
			name:        "Valid database path",
			databaseURL: filepath.Join(t.TempDir(), "test.db"),
			wantErr:     false,
		},
		{
			name:        "Default database path",
			databaseURL: "",
			wantErr:     false,
		},
		{
			name:        "Invalid database path",
			databaseURL: "/invalid/path/that/does/not/exist/test.db",
			wantErr:     true, // SQLite can't create file in non-existent directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			db, err := New(ctx, tt.databaseURL)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if db != nil {
					db.Close()
					// Clean up default database if created
					if tt.databaseURL == "" {
						os.Remove("./cms.db")
					}
				}
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		post    Post
		wantErr bool
	}{
		{
			name: "Valid post with all fields",
			post: Post{
				Title:   "Test Post",
				Content: sql.NullString{String: "Test content", Valid: true},
				Author:  sql.NullString{String: "Test Author", Valid: true},
				Slug:    sql.NullString{String: "test-post", Valid: true},
			},
			wantErr: false,
		},
		{
			name: "Post with only required fields",
			post: Post{
				Title: "Minimal Post",
			},
			wantErr: false,
		},
		{
			name: "Post with empty title",
			post: Post{
				Title: "",
			},
			wantErr: false, // SQLite allows empty strings
		},
		{
			name: "Duplicate slug",
			post: Post{
				Title: "Duplicate Post",
				Slug:  sql.NullString{String: "duplicate-slug", Valid: true},
			},
			wantErr: false, // First insert should succeed
		},
	}

	// Create a post with duplicate slug first
	_, err := db.CreatePost(ctx, Post{
		Title: "First Post",
		Slug:  sql.NullString{String: "duplicate-slug", Valid: true},
	})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if testing duplicate slug on second run
			if tt.name == "Duplicate slug" {
				_, err := db.CreatePost(ctx, tt.post)
				assert.Error(t, err, "Should fail on duplicate slug")
				return
			}

			createdPost, err := db.CreatePost(ctx, tt.post)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, createdPost)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, createdPost)
				
				// Verify created post fields
				assert.Equal(t, tt.post.Title, createdPost.Title)
				assert.Equal(t, tt.post.Content, createdPost.Content)
				assert.Equal(t, tt.post.Author, createdPost.Author)
				assert.Equal(t, tt.post.Slug, createdPost.Slug)
				
				// Verify auto-generated fields
				assert.Greater(t, createdPost.ID, int64(0))
				assert.True(t, createdPost.CreatedAt.Valid)
				assert.True(t, createdPost.UpdatedAt.Valid)
				assert.WithinDuration(t, time.Now(), createdPost.CreatedAt.Time, 2*time.Second)
				assert.WithinDuration(t, time.Now(), createdPost.UpdatedAt.Time, 2*time.Second)
			}
		})
	}
}

func TestGetPostByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test posts
	post1, err := db.CreatePost(ctx, Post{
		Title:   "Test Post 1",
		Content: sql.NullString{String: "Content 1", Valid: true},
		Author:  sql.NullString{String: "Author 1", Valid: true},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid ID",
			id:      int(post1.ID),
			wantErr: false,
		},
		{
			name:    "Non-existent ID",
			id:      9999,
			wantErr: true,
			errMsg:  "post not found",
		},
		{
			name:    "Zero ID",
			id:      0,
			wantErr: true,
			errMsg:  "post not found",
		},
		{
			name:    "Negative ID",
			id:      -1,
			wantErr: true,
			errMsg:  "post not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := db.GetPostByID(ctx, tt.id)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, post)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, post)
				assert.Equal(t, post1.ID, post.ID)
				assert.Equal(t, post1.Title, post.Title)
			}
		})
	}
}

func TestGetPostBySlug(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test post
	post1, err := db.CreatePost(ctx, Post{
		Title: "Test Post",
		Slug:  sql.NullString{String: "test-slug", Valid: true},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		slug    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid slug",
			slug:    "test-slug",
			wantErr: false,
		},
		{
			name:    "Non-existent slug",
			slug:    "non-existent",
			wantErr: true,
			errMsg:  "post not found",
		},
		{
			name:    "Empty slug",
			slug:    "",
			wantErr: true,
			errMsg:  "post not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := db.GetPostBySlug(ctx, tt.slug)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, post)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, post)
				assert.Equal(t, post1.ID, post.ID)
				assert.Equal(t, tt.slug, post.Slug.String)
			}
		})
	}
}

func TestUpdatePostByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create initial post
	initialPost, err := db.CreatePost(ctx, Post{
		Title:   "Original Title",
		Content: sql.NullString{String: "Original Content", Valid: true},
		Author:  sql.NullString{String: "Original Author", Valid: true},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int
		updates Post
		wantErr bool
		errMsg  string
	}{
		{
			name: "Update all fields",
			id:   int(initialPost.ID),
			updates: Post{
				Title:   "Updated Title",
				Content: sql.NullString{String: "Updated Content", Valid: true},
				Author:  sql.NullString{String: "Updated Author", Valid: true},
			},
			wantErr: false,
		},
		{
			name: "Update only title",
			id:   int(initialPost.ID),
			updates: Post{
				Title:   "New Title Only",
				Content: initialPost.Content,
				Author:  initialPost.Author,
			},
			wantErr: false,
		},
		{
			name:    "Update non-existent post",
			id:      9999,
			updates: Post{Title: "Won't work"},
			wantErr: true,
			errMsg:  "post not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedPost, err := db.UpdatePostByID(ctx, tt.id, tt.updates)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, updatedPost)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, updatedPost)
				
				// Verify updates
				assert.Equal(t, tt.updates.Title, updatedPost.Title)
				assert.Equal(t, tt.updates.Content, updatedPost.Content)
				assert.Equal(t, tt.updates.Author, updatedPost.Author)
				
				// Verify UpdatedAt was changed
				assert.True(t, updatedPost.UpdatedAt.Time.After(initialPost.UpdatedAt.Time))
			}
		})
	}
}

func TestDeletePostByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test post
	post, err := db.CreatePost(ctx, Post{
		Title: "To Be Deleted",
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Delete existing post",
			id:      int(post.ID),
			wantErr: false,
		},
		{
			name:    "Delete non-existent post",
			id:      9999,
			wantErr: false, // SQLite doesn't return error for no rows affected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.DeletePostByID(ctx, tt.id)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify post was deleted if it existed
				if tt.id == int(post.ID) {
					_, err := db.GetPostByID(ctx, tt.id)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "post not found")
				}
			}
		})
	}
}

func TestListPosts(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		limit       int
		offset      int
		wantCount   int
		wantErr     bool
	}{
		{
			name:      "List all posts",
			limit:     0,
			offset:    0,
			wantCount: 7, // 5 created + 2 from sample data migration
			wantErr:   false,
		},
		{
			name:      "List with limit",
			limit:     3,
			offset:    0,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "List with limit and offset",
			limit:     2,
			offset:    2,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "List with offset beyond total",
			limit:     10,
			offset:    10,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh database for each test
			db, cleanup := setupTestDB(t)
			defer cleanup()
			
			// Create test posts
			postCount := 5
			for i := 0; i < postCount; i++ {
				_, err := db.CreatePost(ctx, Post{
					Title: fmt.Sprintf("Post %d", i+1),
				})
				require.NoError(t, err)
			}
			
			posts, err := db.ListPosts(ctx, tt.limit, tt.offset)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, posts)
			} else {
				assert.NoError(t, err)
				t.Logf("Expected %d posts, got %d posts", tt.wantCount, len(posts))
				assert.Len(t, posts, tt.wantCount)
				
				// Verify we got posts (ordering tests are complex due to different SQL queries)
				if tt.wantCount > 0 {
					assert.NotEmpty(t, posts, "Should return posts when count > 0")
				}
			}
		})
	}
}

func TestDatabaseClose(t *testing.T) {
	db, _ := setupTestDB(t)
	
	// Close should not error
	err := db.Close()
	assert.NoError(t, err)
	
	// Operations after close should fail
	ctx := context.Background()
	_, err = db.ListPosts(ctx, 0, 0)
	assert.Error(t, err)
}