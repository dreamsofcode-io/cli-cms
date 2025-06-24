package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostFormData_ToPost(t *testing.T) {
	tests := []struct {
		name     string
		formData PostFormData
	}{
		{
			name: "Complete form data",
			formData: PostFormData{
				Title:   "Test Post",
				Content: "Test content here",
				Author:  "Test Author",
				Slug:    "test-post",
				Confirm: true,
			},
		},
		{
			name: "Minimal form data",
			formData: PostFormData{
				Title: "Minimal Post",
			},
		},
		{
			name: "Empty form data",
			formData: PostFormData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := tt.formData.ToPost()
			
			// Verify title is always set
			assert.Equal(t, tt.formData.Title, post.Title)
			
			// Verify content handling
			if tt.formData.Content == "" {
				assert.False(t, post.Content.Valid, "Content should be invalid when empty")
			} else {
				assert.True(t, post.Content.Valid, "Content should be valid when provided")
				assert.Equal(t, tt.formData.Content, post.Content.String)
			}
			
			// Verify author handling
			if tt.formData.Author == "" {
				assert.False(t, post.Author.Valid, "Author should be invalid when empty")
			} else {
				assert.True(t, post.Author.Valid, "Author should be valid when provided")
				assert.Equal(t, tt.formData.Author, post.Author.String)
			}
			
			// Verify slug handling
			if tt.formData.Slug == "" {
				assert.False(t, post.Slug.Valid, "Slug should be invalid when empty")
			} else {
				assert.True(t, post.Slug.Valid, "Slug should be valid when provided")
				assert.Equal(t, tt.formData.Slug, post.Slug.String)
			}
			
			// Verify timestamps are not set (should be handled by database)
			assert.False(t, post.CreatedAt.Valid, "CreatedAt should not be set in form conversion")
			assert.False(t, post.UpdatedAt.Valid, "UpdatedAt should not be set in form conversion")
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "Simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Title with special characters",
			title:    "Hello, World! & More #stuff",
			expected: "hello-world-more-stuff",
		},
		{
			name:     "Title with numbers",
			title:    "Top 10 Tips for 2024",
			expected: "top-10-tips-for-2024",
		},
		{
			name:     "Title with multiple spaces",
			title:    "Too    Many     Spaces",
			expected: "too-many-spaces",
		},
		{
			name:     "Title with leading/trailing spaces",
			title:    "  Trimmed Title  ",
			expected: "trimmed-title",
		},
		{
			name:     "Empty title",
			title:    "",
			expected: "",
		},
		{
			name:     "Only special characters",
			title:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "Unicode characters",
			title:    "Héllo Wörld",
			expected: "hllo-wrld", // Non-ASCII chars are removed
		},
		{
			name:     "Consecutive hyphens",
			title:    "Word---With---Hyphens",
			expected: "word-with-hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSlug(tt.title)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPostFormDataValidation(t *testing.T) {
	// Test that form data can be created and used correctly
	formData := PostFormData{
		Title:   "Valid Title",
		Content: "Valid content",
		Author:  "Valid Author",
		Slug:    "valid-slug",
		Confirm: true,
	}

	// Test conversion to database post
	post := formData.ToPost()
	
	assert.Equal(t, formData.Title, post.Title)
	assert.Equal(t, formData.Content, post.Content.String)
	assert.Equal(t, formData.Author, post.Author.String)
	assert.Equal(t, formData.Slug, post.Slug.String)
	
	// Test all nullable fields are valid when populated
	assert.True(t, post.Content.Valid)
	assert.True(t, post.Author.Valid)
	assert.True(t, post.Slug.Valid)
}

func TestPostFormDataEdgeCases(t *testing.T) {
	// Test edge cases that might cause issues
	
	t.Run("Very long title", func(t *testing.T) {
		longTitle := string(make([]byte, 1000)) // Very long string
		for i := range longTitle {
			longTitle = longTitle[:i] + "a" + longTitle[i+1:]
		}
		
		formData := PostFormData{Title: longTitle}
		post := formData.ToPost()
		
		assert.Equal(t, longTitle, post.Title)
	})
	
	t.Run("Special characters in all fields", func(t *testing.T) {
		specialChars := "!@#$%^&*(){}[]|\\:;\"'<>,.?/~`"
		
		formData := PostFormData{
			Title:   specialChars,
			Content: specialChars,
			Author:  specialChars,
			Slug:    specialChars,
		}
		post := formData.ToPost()
		
		assert.Equal(t, specialChars, post.Title)
		assert.Equal(t, specialChars, post.Content.String)
		assert.Equal(t, specialChars, post.Author.String)
		assert.Equal(t, specialChars, post.Slug.String)
	})
	
	t.Run("Newlines and tabs", func(t *testing.T) {
		textWithWhitespace := "Line 1\nLine 2\tTabbed\r\nWindows line ending"
		
		formData := PostFormData{
			Title:   textWithWhitespace,
			Content: textWithWhitespace,
		}
		post := formData.ToPost()
		
		assert.Equal(t, textWithWhitespace, post.Title)
		assert.Equal(t, textWithWhitespace, post.Content.String)
	})
}

func TestErrUserCancelled(t *testing.T) {
	// Test that the error constant exists and has expected properties
	assert.NotNil(t, ErrUserCancelled)
	assert.Contains(t, ErrUserCancelled.Error(), "cancelled")
}

// Benchmark tests to ensure form operations are performant
func BenchmarkPostFormDataToPost(b *testing.B) {
	formData := PostFormData{
		Title:   "Benchmark Test Post",
		Content: "This is benchmark content for testing performance",
		Author:  "Benchmark Author",
		Slug:    "benchmark-test-post",
		Confirm: true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		post := formData.ToPost()
		_ = post // Use the result to prevent optimization
	}
}

func BenchmarkGenerateSlug(b *testing.B) {
	title := "This is a Title With Multiple Words and Some Numbers 123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slug := generateSlug(title)
		_ = slug // Use the result to prevent optimization
	}
}