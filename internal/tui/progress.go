package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressModel represents the state of our progress bar TUI application
type ProgressModel struct {
	progress   progress.Model
	isQuitting bool
	message    string
	percent    float64
}

// TickMsg represents a tick for progress updates
type TickMsg time.Time

// InitialProgressModel creates a new progress model with default configuration
func InitialProgressModel(message string) ProgressModel {
	return ProgressModel{
		progress: progress.New(progress.WithDefaultGradient()),
		message:  message,
		percent:  0.0,
	}
}

// Init implements the tea.Model interface - returns initial command
func (m ProgressModel) Init() tea.Cmd {
	return tickCmd()
}

// Update implements the tea.Model interface - handles messages and updates state
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4 // Account for padding
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil
	
	case TickMsg:
		if m.percent >= 1.0 {
			m.isQuitting = true
			return m, tea.Quit
		}
		
		// Increment progress by 2% each tick
		m.percent += 0.02
		if m.percent > 1.0 {
			m.percent = 1.0
		}
		
		return m, tickCmd()
	
	default:
		return m, nil
	}
}

// View implements the tea.Model interface - renders the current state
func (m ProgressModel) View() string {
	if m.isQuitting && m.percent >= 1.0 {
		return fmt.Sprintf("✅ %s completed!\n", m.message)
	}
	
	if m.isQuitting {
		return fmt.Sprintf("❌ %s cancelled\n", m.message)
	}
	
	pad := lipgloss.NewStyle().Padding(1, 2)
	
	return pad.Render(fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.message,
		m.progress.ViewAs(m.percent),
		fmt.Sprintf("%.0f%% complete (Press q to cancel)", m.percent*100),
	))
}

// tickCmd returns a command that sends a TickMsg after a delay
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// RunProgress starts a progress bar with the given message
func RunProgress(message string) error {
	model := InitialProgressModel(message)
	p := tea.NewProgram(model)
	
	_, err := p.Run()
	return err
}