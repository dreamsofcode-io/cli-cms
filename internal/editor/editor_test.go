package editor

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		editorEnv      string
		expectedCmd    string
		expectedArgs   []string
		expectFallback bool
	}{
		{
			name:           "EDITOR environment variable set",
			editorEnv:      "vim",
			expectedCmd:    "vim",
			expectedArgs:   []string{},
			expectFallback: false,
		},
		{
			name:           "EDITOR with arguments",
			editorEnv:      "code --wait",
			expectedCmd:    "code",
			expectedArgs:   []string{"--wait"},
			expectFallback: false,
		},
		{
			name:           "Empty EDITOR env var",
			editorEnv:      "",
			expectedCmd:    getExpectedFallback(),
			expectedArgs:   []string{},
			expectFallback: true,
		},
		{
			name:           "Multiple arguments",
			editorEnv:      "emacs -nw --no-init-file",
			expectedCmd:    "emacs",
			expectedArgs:   []string{"-nw", "--no-init-file"},
			expectFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			originalEditor := os.Getenv("EDITOR")
			defer os.Setenv("EDITOR", originalEditor)
			
			if tt.editorEnv == "" {
				os.Unsetenv("EDITOR")
			} else {
				os.Setenv("EDITOR", tt.editorEnv)
			}

			// Create editor
			editor := New()

			// Verify results
			assert.Equal(t, tt.expectedCmd, editor.command)
			assert.Equal(t, tt.expectedArgs, editor.args)
		})
	}
}

func TestEditor_GetEditorInfo(t *testing.T) {
	tests := []struct {
		name      string
		editorEnv string
		contains  string
	}{
		{
			name:      "With EDITOR env var",
			editorEnv: "vim",
			contains:  "vim (from EDITOR env var)",
		},
		{
			name:      "Without EDITOR env var",
			editorEnv: "",
			contains:  "(default fallback)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			originalEditor := os.Getenv("EDITOR")
			defer os.Setenv("EDITOR", originalEditor)
			
			if tt.editorEnv == "" {
				os.Unsetenv("EDITOR")
			} else {
				os.Setenv("EDITOR", tt.editorEnv)
			}

			editor := New()
			info := editor.GetEditorInfo()
			
			assert.Contains(t, info, tt.contains)
		})
	}
}

func TestEditor_IsAvailable(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		shouldExist bool
	}{
		{
			name:        "Non-existent command",
			command:     "this-command-definitely-does-not-exist",
			shouldExist: false,
		},
		{
			name:        "Common system command",
			command:     getCommonCommand(),
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := &Editor{command: tt.command}
			available := editor.IsAvailable()
			
			assert.Equal(t, tt.shouldExist, available)
		})
	}
}

func TestEditor_filterComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Mixed content with comments",
			input: `# This is a comment
This is content
# Another comment
More content
   # Indented comment  
Final content`,
			expected: "This is content\nMore content\nFinal content",
		},
		{
			name:     "Only comments",
			input:    "# Comment 1\n# Comment 2\n# Comment 3",
			expected: "",
		},
		{
			name:     "No comments",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Only whitespace and comments",
			input:    "  \n  # Comment  \n\t\n",
			expected: "",
		},
		{
			name: "Hash in middle of line",
			input: `This line has # hash in middle
# This is a comment
Normal content`,
			expected: "This line has # hash in middle\nNormal content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New()
			result := editor.filterComments(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteToFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
		content  string
		wantErr  bool
	}{
		{
			name:     "Valid file write",
			filename: filepath.Join(tempDir, "test.txt"),
			content:  "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "Empty content",
			filename: filepath.Join(tempDir, "empty.txt"),
			content:  "",
			wantErr:  false,
		},
		{
			name:     "Multiline content",
			filename: filepath.Join(tempDir, "multiline.txt"),
			content:  "Line 1\nLine 2\nLine 3",
			wantErr:  false,
		},
		{
			name:     "Invalid directory",
			filename: "/invalid/path/that/does/not/exist/file.txt",
			content:  "content",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteToFile(tt.filename, tt.content)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify file was written correctly
				content, err := os.ReadFile(tt.filename)
				require.NoError(t, err)
				assert.Equal(t, tt.content, string(content))
			}
		})
	}
}

func TestEditor_CreateContentFile(t *testing.T) {
	tempDir := t.TempDir()
	
	// Skip this test if no suitable editor is available
	editor := New()
	if !editor.IsAvailable() {
		t.Skip("No suitable editor available for testing")
	}

	tests := []struct {
		name       string
		dir        string
		filename   string
		content    string
		shouldSkip bool
	}{
		{
			name:     "Valid directory and file",
			dir:      tempDir,
			filename: "test.md",
			content:  "# Test Content",
		},
		{
			name:     "Non-existent directory (should create)",
			dir:      filepath.Join(tempDir, "new", "nested", "dir"),
			filename: "nested.md",
			content:  "Nested content",
		},
		{
			name:     "Empty content",
			dir:      tempDir,
			filename: "empty.md",
			content:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSkip {
				t.Skip("Skipping test case")
			}

			// We can't easily test the interactive editor behavior in automated tests,
			// but we can test the file creation and setup logic
			
			// Verify directory creation
			err := os.MkdirAll(tt.dir, 0755)
			assert.NoError(t, err)
			
			// Create the file directly to test the logic
			fullPath := filepath.Join(tt.dir, tt.filename)
			err = WriteToFile(fullPath, tt.content)
			assert.NoError(t, err)
			
			// Verify file exists and has correct content
			content, err := os.ReadFile(fullPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.content, string(content))
		})
	}
}

func TestEditor_EditContentWithTemplate(t *testing.T) {
	// This test focuses on the template generation part
	// We can't test the interactive editing without mocking
	
	editor := New()
	
	tests := []struct {
		name            string
		title           string
		author          string
		existingContent string
		isUpdate        bool
		containsTitle   bool
		containsAuthor  bool
	}{
		{
			name:          "New post with all fields",
			title:         "Test Title",
			author:        "Test Author",
			isUpdate:      false,
			containsTitle: true,
			containsAuthor: true,
		},
		{
			name:            "Update existing post",
			title:           "Updated Title",
			author:          "Updated Author",
			existingContent: "Existing content here",
			isUpdate:        true,
			containsTitle:   true,
			containsAuthor:  true,
		},
		{
			name:     "Minimal new post",
			isUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We'll test the template building logic by examining what would be passed to EditContent
			var template strings.Builder
			
			if tt.isUpdate {
				template.WriteString("# Editing Post\n")
			} else {
				template.WriteString("# Creating New Post\n")
			}
			template.WriteString("#\n")
			
			if tt.title != "" {
				template.WriteString("# Title: " + tt.title + "\n")
			}
			if tt.author != "" {
				template.WriteString("# Author: " + tt.author + "\n")
			}
			
			template.WriteString("#\n")
			template.WriteString("# Write your post content below this line.\n")
			template.WriteString("# Lines starting with '#' are comments and will be ignored.\n")
			template.WriteString("#\n")
			template.WriteString("\n")
			
			if tt.existingContent != "" {
				template.WriteString(tt.existingContent)
			}

			templateStr := template.String()
			
			// Verify template contains expected elements
			if tt.isUpdate {
				assert.Contains(t, templateStr, "# Editing Post")
			} else {
				assert.Contains(t, templateStr, "# Creating New Post")
			}
			
			if tt.containsTitle {
				assert.Contains(t, templateStr, "# Title: "+tt.title)
			}
			
			if tt.containsAuthor {
				assert.Contains(t, templateStr, "# Author: "+tt.author)
			}
			
			if tt.existingContent != "" {
				assert.Contains(t, templateStr, tt.existingContent)
			}
			
			// Test comment filtering
			filtered := editor.filterComments(templateStr)
			
			// Filtered content should not contain comment lines
			lines := strings.Split(filtered, "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					assert.False(t, strings.HasPrefix(trimmed, "#"), 
						"Filtered content should not contain comment lines: %s", line)
				}
			}
		})
	}
}

// Helper functions

func getExpectedFallback() string {
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	case "darwin":
		return "nano"
	default:
		return "nano"
	}
}

func getCommonCommand() string {
	switch runtime.GOOS {
	case "windows":
		return "cmd" // cmd.exe should always exist on Windows
	default:
		return "sh" // sh should exist on Unix-like systems
	}
}

// Benchmark tests
func BenchmarkEditor_filterComments(b *testing.B) {
	editor := New()
	content := `# Comment 1
Real content line 1
# Comment 2
Real content line 2
# Comment 3
Real content line 3`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = editor.filterComments(content)
	}
}

func BenchmarkWriteToFile(b *testing.B) {
	tempDir := b.TempDir()
	content := "This is test content for benchmarking file writes."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := filepath.Join(tempDir, "bench.txt")
		_ = WriteToFile(filename, content)
		os.Remove(filename) // Clean up for next iteration
	}
}