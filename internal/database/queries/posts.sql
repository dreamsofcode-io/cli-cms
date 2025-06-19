-- name: ListPosts :many
SELECT * FROM posts ORDER BY id ASC;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = ?;

-- name: GetPostBySlug :one
SELECT * FROM posts WHERE slug = ?;

-- name: CreatePost :one
INSERT INTO posts (title, content, author, slug, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePostByID :one
UPDATE posts 
SET title = ?, content = ?, author = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdatePostBySlug :one
UPDATE posts 
SET title = ?, content = ?, author = ?, updated_at = ?
WHERE slug = ?
RETURNING *;

-- name: DeletePostByID :exec
DELETE FROM posts WHERE id = ?;

-- name: DeletePostBySlug :exec
DELETE FROM posts WHERE slug = ?;

-- name: ListPostsWithPagination :many
SELECT * FROM posts 
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;