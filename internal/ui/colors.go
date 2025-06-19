package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Color functions for consistent styling across the application
var (
	// Success messages (green)
	Success = color.New(color.FgGreen, color.Bold).PrintfFunc()
	SuccessString = color.New(color.FgGreen, color.Bold).SprintfFunc()
	
	// Error messages (red)
	Error = color.New(color.FgRed, color.Bold).PrintfFunc()
	ErrorString = color.New(color.FgRed, color.Bold).SprintfFunc()
	
	// Warning messages (yellow)
	Warning = color.New(color.FgYellow, color.Bold).PrintfFunc()
	WarningString = color.New(color.FgYellow, color.Bold).SprintfFunc()
	
	// Info messages (blue)
	Info = color.New(color.FgBlue).PrintfFunc()
	InfoString = color.New(color.FgBlue).SprintfFunc()
	
	// Highlight/emphasis (cyan)
	Highlight = color.New(color.FgCyan, color.Bold).PrintfFunc()
	HighlightString = color.New(color.FgCyan, color.Bold).SprintfFunc()
	
	// Subtle/secondary text (dim white)
	Subtle = color.New(color.FgWhite, color.Faint).PrintfFunc()
	SubtleString = color.New(color.FgWhite, color.Faint).SprintfFunc()
	
	// Links (blue with underline)
	Link = color.New(color.FgBlue, color.Underline).PrintfFunc()
	LinkString = color.New(color.FgBlue, color.Underline).SprintfFunc()
)

// Helper functions for common output patterns

// PrintSuccess prints a success message with checkmark
func PrintSuccess(format string, args ...interface{}) {
	Success("✅ "+format, args...)
}

// PrintError prints an error message with X mark
func PrintError(format string, args ...interface{}) {
	Error("❌ "+format, args...)
}

// PrintWarning prints a warning message with warning sign
func PrintWarning(format string, args ...interface{}) {
	Warning("⚠️  "+format, args...)
}

// PrintInfo prints an info message with info icon
func PrintInfo(format string, args ...interface{}) {
	Info("ℹ️  "+format, args...)
}

// ErrorAndExit prints an error message and exits with code 1
func ErrorAndExit(format string, args ...interface{}) {
	PrintError(format, args...)
	os.Exit(1)
}

// Header prints a section header
func Header(text string) {
	fmt.Printf("\n")
	Highlight("=== %s ===\n", text)
	fmt.Printf("\n")
}

// Subheader prints a subsection header
func Subheader(text string) {
	fmt.Printf("\n")
	Info("--- %s ---\n", text)
}

// Field prints a label-value pair with consistent formatting
func Field(label string, value interface{}) {
	fmt.Printf("%s: %v\n", HighlightString(label), value)
}

// FieldIfNotEmpty prints a field only if the value is not empty
func FieldIfNotEmpty(label, value string) {
	if value != "" {
		Field(label, value)
	}
}