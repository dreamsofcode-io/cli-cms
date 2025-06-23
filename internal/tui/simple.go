package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// DemoSpinner shows what a spinner would look like without full TUI
func DemoSpinner(message string, duration time.Duration) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red

	fmt.Printf("🎭 TUI Spinner Demo\n")
	fmt.Printf("Message: %s\n", message)
	fmt.Printf("Spinner frames: %v\n", s.Spinner.Frames)
	fmt.Printf("Spinner FPS: %v\n", s.Spinner.FPS)
	
	// Simulate a few spinner frames
	for i := 0; i < 10; i++ {
		frame := s.Spinner.Frames[i%len(s.Spinner.Frames)]
		styled := s.Style.Render(frame)
		fmt.Printf("Frame %d: %s %s\n", i+1, styled, message)
		time.Sleep(time.Duration(1000/s.Spinner.FPS) * time.Millisecond)
	}
	
	fmt.Printf("✅ Spinner demo completed!\n")
}

// DemoProgress shows what a progress bar would look like without full TUI
func DemoProgress(message string) {
	fmt.Printf("📊 TUI Progress Demo\n")
	fmt.Printf("Message: %s\n", message)
	
	// Simulate progress updates
	for i := 0; i <= 10; i++ {
		percent := float64(i) / 10.0
		
		// Create a simple ASCII progress bar
		filled := int(percent * 20)
		bar := ""
		for j := 0; j < 20; j++ {
			if j < filled {
				bar += "█"
			} else {
				bar += "░"
			}
		}
		
		fmt.Printf("\r[%s] %.0f%% %s", bar, percent*100, message)
		
		if i < 10 {
			time.Sleep(200 * time.Millisecond)
		}
	}
	
	fmt.Printf("\n✅ Progress demo completed!\n")
}

// DemoCounter shows what a counter would look like without full TUI
func DemoCounter(duration time.Duration) {
	fmt.Printf("🔢 TUI Counter Demo\n")
	fmt.Printf("Duration: %v\n", duration)
	
	start := time.Now()
	count := 0
	
	for time.Since(start) < duration {
		count++
		
		// Style the counter with lipgloss
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")). // Bright blue
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")) // Gray border
		
		content := fmt.Sprintf("Count: %d", count)
		styled := style.Render(content)
		
		fmt.Printf("\r%s", styled)
		time.Sleep(time.Second)
	}
	
	fmt.Printf("\n✅ Counter demo completed!\n")
}