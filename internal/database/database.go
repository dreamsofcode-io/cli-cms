package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Post represents a blog post in our CMS
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Database wraps the sql.DB connection and provides methods for database operations
type Database struct {
	db *sql.DB
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

	database := &Database{db: db}

	// Initialize the schema
	if err := database.createTables(ctx); err != nil {
		return nil, err
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// createTables creates the posts table if it doesn't exist
func (d *Database) createTables(ctx context.Context) error {
	const createPostsTable = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT,
		author TEXT,
		slug TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.ExecContext(ctx, createPostsTable)
	return err
}

// CreatePost inserts a new post into the database
func (d *Database) CreatePost(ctx context.Context, post Post) (*Post, error) {
	const insertPost = `
	INSERT INTO posts (title, content, author, slug, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	result, err := d.db.ExecContext(ctx, insertPost, 
		post.Title, post.Content, post.Author, post.Slug, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	post.ID = int(id)
	return &post, nil
}

// GetPostByID retrieves a post by its ID
func (d *Database) GetPostByID(ctx context.Context, id int) (*Post, error) {
	const selectPost = `
	SELECT id, title, content, author, slug, created_at, updated_at 
	FROM posts WHERE id = ?`

	row := d.db.QueryRowContext(ctx, selectPost, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, 
		&post.Slug, &post.CreatedAt, &post.UpdatedAt)
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
	const selectPost = `
	SELECT id, title, content, author, slug, created_at, updated_at 
	FROM posts WHERE slug = ?`

	row := d.db.QueryRowContext(ctx, selectPost, slug)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, 
		&post.Slug, &post.CreatedAt, &post.UpdatedAt)
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
	updates.UpdatedAt = time.Now()
	
	const updatePost = `
	UPDATE posts 
	SET title = ?, content = ?, author = ?, updated_at = ?
	WHERE id = ?`

	result, err := d.db.ExecContext(ctx, updatePost, 
		updates.Title, updates.Content, updates.Author, updates.UpdatedAt, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("post not found")
	}

	// Return the updated post
	return d.GetPostByID(ctx, id)
}

// UpdatePostBySlug updates a post by its slug
func (d *Database) UpdatePostBySlug(ctx context.Context, slug string, updates Post) (*Post, error) {
	// First get the post to find its ID
	existingPost, err := d.GetPostBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return d.UpdatePostByID(ctx, existingPost.ID, updates)
}

// DeletePostByID deletes a post by its ID
func (d *Database) DeletePostByID(ctx context.Context, id int) error {
	const deletePost = `DELETE FROM posts WHERE id = ?`

	result, err := d.db.ExecContext(ctx, deletePost, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}

// DeletePostBySlug deletes a post by its slug
func (d *Database) DeletePostBySlug(ctx context.Context, slug string) error {
	const deletePost = `DELETE FROM posts WHERE slug = ?`

	result, err := d.db.ExecContext(ctx, deletePost, slug)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}

// ListPosts retrieves all posts with optional limit and offset for pagination
func (d *Database) ListPosts(ctx context.Context, limit, offset int) ([]*Post, error) {
	const selectAllPosts = `
	SELECT id, title, content, author, slug, created_at, updated_at 
	FROM posts 
	ORDER BY created_at DESC 
	LIMIT ? OFFSET ?`

	rows, err := d.db.QueryContext(ctx, selectAllPosts, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Always Be Closing

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, 
			&post.Slug, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}