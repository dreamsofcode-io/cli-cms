package database

import (
	"database/sql"
	"time"
)

// Helper functions to convert between string and sql.NullString
func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func NullStringToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// Helper functions to convert between time.Time and sql.NullTime
func TimeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func NullTimeToTime(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

// CreatePostFromInput creates a Post from simple input parameters
func CreatePostFromInput(title, content, author, slug string) Post {
	return Post{
		Title:   title,
		Content: StringToNullString(content),
		Author:  StringToNullString(author),
		Slug:    StringToNullString(slug),
	}
}