-- name: ListPosts :many
SELECT * FROM posts ORDER BY id DESC;

-- name: InsertPost :exec
INSERT INTO posts (title, slug, content) VALUES (
  ?, ?, ?
);
