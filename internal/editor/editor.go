package editor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Editor handles opening text editors for content editing
type Editor struct {
	command string
	args    []string
}

// New creates a new Editor instance using the EDITOR environment variable
// Falls back to common editors if EDITOR is not set
func New() *Editor {
	editorCmd := os.Getenv("EDITOR")
	
	if editorCmd == "" {
		// Fallback to common editors based on OS
		switch runtime.GOOS {
		case "windows":
			editorCmd = "notepad"
		case "darwin":
			editorCmd = "nano"
		default:
			editorCmd = "nano"
		}
	}

	// Parse command and arguments
	parts := strings.Fields(editorCmd)
	if len(parts) == 0 {
		parts = []string{"nano"} // ultimate fallback
	}

	return &Editor{
		command: parts[0],
		args:    parts[1:],
	}
}

// EditContent opens an editor with initial content and returns the edited content
func (e *Editor) EditContent(initialContent string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "cms-edit-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up temp file

	// Write initial content to temp file
	if initialContent != "" {
		if _, err := tmpFile.WriteString(initialContent); err != nil {
			tmpFile.Close()
			return "", fmt.Errorf("failed to write initial content: %w", err)
		}
	}

	// Close the file so the editor can open it
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	// Prepare editor command with temp file
	args := append(e.args, tmpFile.Name())
	cmd := exec.Command(e.command, args...)
	
	// Connect editor to terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the editor
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the edited content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited content: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

// EditContentWithTemplate opens an editor with a template and returns the edited content
func (e *Editor) EditContentWithTemplate(title, author, existingContent string, isUpdate bool) (string, error) {
	var template strings.Builder
	
	// Add template content
	if isUpdate {
		template.WriteString("# Editing Post\n")
	} else {
		template.WriteString("# Creating New Post\n")
	}
	template.WriteString("#\n")
	
	if title != "" {
		template.WriteString(fmt.Sprintf("# Title: %s\n", title))
	}
	if author != "" {
		template.WriteString(fmt.Sprintf("# Author: %s\n", author))
	}
	
	template.WriteString("#\n")
	template.WriteString("# Write your post content below this line.\n")
	template.WriteString("# Lines starting with '#' are comments and will be ignored.\n")
	template.WriteString("#\n")
	template.WriteString("\n")
	
	// Add existing content if provided
	if existingContent != "" {
		template.WriteString(existingContent)
	}

	// Edit the template
	editedContent, err := e.EditContent(template.String())
	if err != nil {
		return "", err
	}

	// Filter out comment lines and return clean content
	return e.filterComments(editedContent), nil
}

// filterComments removes lines starting with '#' and trims whitespace
func (e *Editor) filterComments(content string) string {
	lines := strings.Split(content, "\n")
	var filtered []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and comment lines
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			filtered = append(filtered, line)
		}
	}
	
	return strings.TrimSpace(strings.Join(filtered, "\n"))
}

// GetEditorInfo returns information about the configured editor
func (e *Editor) GetEditorInfo() string {
	editorEnv := os.Getenv("EDITOR")
	if editorEnv != "" {
		return fmt.Sprintf("%s (from EDITOR env var)", editorEnv)
	}
	return fmt.Sprintf("%s (default fallback)", e.command)
}

// IsAvailable checks if the configured editor is available on the system
func (e *Editor) IsAvailable() bool {
	_, err := exec.LookPath(e.command)
	return err == nil
}

// CreateContentFile creates a content file and opens it in the editor
func (e *Editor) CreateContentFile(dir, filename, content string) (string, error) {
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create full file path
	filePath := filepath.Join(dir, filename)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	// Write initial content
	if content != "" {
		if _, err := file.WriteString(content); err != nil {
			file.Close()
			return "", fmt.Errorf("failed to write content: %w", err)
		}
	}

	file.Close()

	// Open in editor
	args := append(e.args, filePath)
	cmd := exec.Command(e.command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the edited content
	editedContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %w", err)
	}

	return strings.TrimSpace(string(editedContent)), nil
}

// WriteToFile writes content to a file
func WriteToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	return err
}