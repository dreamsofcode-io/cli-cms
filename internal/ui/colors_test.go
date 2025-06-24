package ui

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

// Helper function to capture output using a simpler approach
func captureOutput(fn func()) string {
	// Since the color functions use their own output mechanism,
	// let's use a different approach by temporarily redirecting output
	originalStdout := os.Stdout
	
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_output")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())
	
	// Redirect stdout
	os.Stdout = tmpFile
	
	// Execute function
	fn()
	
	// Restore stdout
	os.Stdout = originalStdout
	tmpFile.Close()
	
	// Read captured output
	output, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		panic(err)
	}
	
	return string(output)
}

func TestColorFunctions(t *testing.T) {
	// Test that the print functions don't panic when called
	// We can't easily test output since fatih/color uses complex output mechanisms
	
	tests := []struct {
		name     string
		function func()
	}{
		{
			name: "PrintSuccess",
			function: func() {
				PrintSuccess("Test success message")
			},
		},
		{
			name: "PrintError",
			function: func() {
				PrintError("Test error message")
			},
		},
		{
			name: "PrintWarning",
			function: func() {
				PrintWarning("Test warning message")
			},
		},
		{
			name: "PrintInfo",
			function: func() {
				PrintInfo("Test info message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that these don't panic
			assert.NotPanics(t, tt.function)
		})
	}
}

func TestPrintFunctionsWithFormatting(t *testing.T) {
	// Test that the print functions with formatting don't panic
	
	tests := []struct {
		name     string
		function func()
	}{
		{
			name: "PrintSuccess with formatting",
			function: func() {
				PrintSuccess("Created %d posts successfully", 5)
			},
		},
		{
			name: "PrintError with formatting", 
			function: func() {
				PrintError("Failed to connect to %s", "database")
			},
		},
		{
			name: "PrintWarning with formatting",
			function: func() {
				PrintWarning("%s is deprecated", "OldFeature")
			},
		},
		{
			name: "PrintInfo with formatting",
			function: func() {
				PrintInfo("Processing %d/%d items", 3, 10)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, tt.function)
		})
	}
}

func TestHeader(t *testing.T) {
	// Test that Header function doesn't panic
	assert.NotPanics(t, func() {
		Header("Test Section")
	})
}

func TestSubheader(t *testing.T) {
	// Test that Subheader function doesn't panic
	assert.NotPanics(t, func() {
		Subheader("Test Subsection")
	})
}

func TestField(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		label    string
		value    interface{}
		expected string
	}{
		{
			name:     "String value",
			label:    "Name",
			value:    "Test Name",
			expected: "Name: Test Name",
		},
		{
			name:     "Integer value",
			label:    "Count",
			value:    42,
			expected: "Count: 42",
		},
		{
			name:     "Boolean value",
			label:    "Active",
			value:    true,
			expected: "Active: true",
		},
		{
			name:     "Nil value",
			label:    "Optional",
			value:    nil,
			expected: "Optional: <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				Field(tt.label, tt.value)
			})
			assert.Contains(t, output, tt.expected)
		})
	}
}

func TestFieldIfNotEmpty(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name          string
		label         string
		value         string
		shouldContain bool
	}{
		{
			name:          "Non-empty value",
			label:         "Description",
			value:         "Test description",
			shouldContain: true,
		},
		{
			name:          "Empty value",
			label:         "Optional",
			value:         "",
			shouldContain: false,
		},
		{
			name:          "Whitespace only",
			label:         "Spaces",
			value:         "   ",
			shouldContain: true, // Function doesn't trim whitespace
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				FieldIfNotEmpty(tt.label, tt.value)
			})

			if tt.shouldContain {
				assert.Contains(t, output, tt.label)
				assert.Contains(t, output, tt.value)
			} else {
				assert.Empty(t, strings.TrimSpace(output))
			}
		})
	}
}

func TestStringFunctions(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		function func() string
		input    string
		expected string
	}{
		{
			name: "SuccessString",
			function: func() string {
				return SuccessString("Success: %s", "completed")
			},
			expected: "Success: completed",
		},
		{
			name: "ErrorString",
			function: func() string {
				return ErrorString("Error: %s", "failed")
			},
			expected: "Error: failed",
		},
		{
			name: "WarningString",
			function: func() string {
				return WarningString("Warning: %s", "deprecated")
			},
			expected: "Warning: deprecated",
		},
		{
			name: "InfoString",
			function: func() string {
				return InfoString("Info: %s", "processing")
			},
			expected: "Info: processing",
		},
		{
			name: "HighlightString",
			function: func() string {
				return HighlightString("Highlighted: %s", "important")
			},
			expected: "Highlighted: important",
		},
		{
			name: "SubtleString",
			function: func() string {
				return SubtleString("Subtle: %s", "note")
			},
			expected: "Subtle: note",
		},
		{
			name: "LinkString",
			function: func() string {
				return LinkString("Link: %s", "https://example.com")
			},
			expected: "Link: https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestColorVariablesExist(t *testing.T) {
	// Test that all color function variables are properly initialized
	assert.NotNil(t, Success)
	assert.NotNil(t, SuccessString)
	assert.NotNil(t, Error)
	assert.NotNil(t, ErrorString)
	assert.NotNil(t, Warning)
	assert.NotNil(t, WarningString)
	assert.NotNil(t, Info)
	assert.NotNil(t, InfoString)
	assert.NotNil(t, Highlight)
	assert.NotNil(t, HighlightString)
	assert.NotNil(t, Subtle)
	assert.NotNil(t, SubtleString)
	assert.NotNil(t, Link)
	assert.NotNil(t, LinkString)
}

func TestColorFunctionsWithColors(t *testing.T) {
	// Test with colors enabled to ensure functions work in color mode
	color.NoColor = false

	tests := []struct {
		name     string
		function func() string
		text     string
	}{
		{
			name: "SuccessString with colors",
			function: func() string {
				return SuccessString("test")
			},
		},
		{
			name: "ErrorString with colors",
			function: func() string {
				return ErrorString("test")
			},
		},
		{
			name: "HighlightString with colors",
			function: func() string {
				return HighlightString("test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			// With colors enabled, the result should contain ANSI escape codes
			// We can't test exact codes as they may vary, but we can test that
			// the function doesn't panic and returns a non-empty string
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "test")
		})
	}
}

func TestEmptyFormatStrings(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		function func()
	}{
		{
			name: "PrintSuccess with empty format",
			function: func() {
				PrintSuccess("")
			},
		},
		{
			name: "PrintError with empty format",
			function: func() {
				PrintError("")
			},
		},
		{
			name: "PrintWarning with empty format",
			function: func() {
				PrintWarning("")
			},
		},
		{
			name: "PrintInfo with empty format",
			function: func() {
				PrintInfo("")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These should not panic
			assert.NotPanics(t, tt.function)
		})
	}
}

// Benchmark tests to ensure color functions are performant
func BenchmarkPrintSuccess(b *testing.B) {
	// Redirect output to avoid cluttering benchmark results
	originalStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = originalStdout }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PrintSuccess("Benchmark test message %d", i)
	}
}

func BenchmarkSuccessString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SuccessString("Benchmark test message %d", i)
	}
}

func BenchmarkField(b *testing.B) {
	// Redirect output to avoid cluttering benchmark results
	originalStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = originalStdout }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Field("Label", fmt.Sprintf("Value %d", i))
	}
}