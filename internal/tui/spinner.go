package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel represents the state of our spinner TUI application
type SpinnerModel struct {
	spinner    spinner.Model
	isQuitting bool
	message    string
}

// InitialSpinnerModel creates a new spinner model with default configuration
func InitialSpinnerModel(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red spinner
	
	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

// Init implements the tea.Model interface - returns initial command
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements the tea.Model interface - handles messages and updates state
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// View implements the tea.Model interface - renders the current state
func (m SpinnerModel) View() string {
	str := fmt.Sprintf("%s %s (Press q to stop)", m.spinner.View(), m.message)
	if m.isQuitting {
		return str + "\n"
	}
	return str
}

// RunSpinner starts a spinner with the given message
func RunSpinner(message string) error {
	model := InitialSpinnerModel(message)
	p := tea.NewProgram(model)
	
	_, err := p.Run()
	return err
}