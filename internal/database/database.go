package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dreamsofcode-io/cli-cms/internal/repository"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Post is an alias for the generated repository Post type
type Post = repository.Post

// Database wraps the sql.DB connection and provides methods for database operations
type Database struct {
	db   *sql.DB
	repo *repository.Queries
}

// New creates a new database connection and initializes the schema
func New(ctx context.Context, databaseURL string) (*Database, error) {
	if databaseURL == "" {
		databaseURL = "./cms.db" // Default SQLite database file
	}

	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Run database migrations
	if err := performMigrations(databaseURL); err != nil {
		return nil, err
	}

	// Create repository
	repo := repository.New(db)

	database := &Database{
		db:   db,
		repo: repo,
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}


// CreatePost inserts a new post into the database
func (d *Database) CreatePost(ctx context.Context, post Post) (*Post, error) {
	now := time.Now()
	
	params := repository.CreatePostParams{
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.Author,
		Slug:      post.Slug,
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}
	
	createdPost, err := d.repo.CreatePost(ctx, params)
	if err != nil {
		return nil, err
	}
	
	return &createdPost, nil
}

// GetPostByID retrieves a post by its ID
func (d *Database) GetPostByID(ctx context.Context, id int) (*Post, error) {
	post, err := d.repo.GetPostByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	
	return &post, nil
}

// GetPostBySlug retrieves a post by its slug
func (d *Database) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	post, err := d.repo.GetPostBySlug(ctx, sql.NullString{String: slug, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	
	return &post, nil
}

// UpdatePostByID updates a post by its ID
func (d *Database) UpdatePostByID(ctx context.Context, id int, updates Post) (*Post, error) {
	now := time.Now()
	
	params := repository.UpdatePostByIDParams{
		ID:        int64(id),
		Title:     updates.Title,
		Content:   updates.Content,
		Author:    updates.Author,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}
	
	updatedPost, err := d.repo.UpdatePostByID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	
	return &updatedPost, nil
}

// UpdatePostBySlug updates a post by its slug
func (d *Database) UpdatePostBySlug(ctx context.Context, slug string, updates Post) (*Post, error) {
	now := time.Now()
	
	params := repository.UpdatePostBySlugParams{
		Slug:      sql.NullString{String: slug, Valid: true},
		Title:     updates.Title,
		Content:   updates.Content,
		Author:    updates.Author,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}
	
	updatedPost, err := d.repo.UpdatePostBySlug(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	
	return &updatedPost, nil
}

// DeletePostByID deletes a post by its ID
func (d *Database) DeletePostByID(ctx context.Context, id int) error {
	err := d.repo.DeletePostByID(ctx, int64(id))
	if err != nil {
		return err
	}
	
	// Note: SQLc doesn't automatically check rows affected for :exec queries
	// You might want to use :execrows if you need to check affected rows
	return nil
}

// DeletePostBySlug deletes a post by its slug
func (d *Database) DeletePostBySlug(ctx context.Context, slug string) error {
	err := d.repo.DeletePostBySlug(ctx, sql.NullString{String: slug, Valid: true})
	if err != nil {
		return err
	}
	
	// Note: SQLc doesn't automatically check rows affected for :exec queries
	// You might want to use :execrows if you need to check affected rows
	return nil
}

// ListPosts retrieves all posts with optional limit and offset for pagination
func (d *Database) ListPosts(ctx context.Context, limit, offset int) ([]*Post, error) {
	if limit > 0 {
		// Use pagination query
		params := repository.ListPostsWithPaginationParams{
			Limit:  int64(limit),
			Offset: int64(offset),
		}
		posts, err := d.repo.ListPostsWithPagination(ctx, params)
		if err != nil {
			return nil, err
		}
		
		// Convert to pointer slice
		result := make([]*Post, len(posts))
		for i := range posts {
			result[i] = &posts[i]
		}
		return result, nil
	}
	
	// Use simple list query (no pagination)
	posts, err := d.repo.ListPosts(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to pointer slice
	result := make([]*Post, len(posts))
	for i := range posts {
		result[i] = &posts[i]
	}
	return result, nil
}