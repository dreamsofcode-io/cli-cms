package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CounterModel represents the state of our counter TUI application
type CounterModel struct {
	count      int
	isQuitting bool
}

// TickCounterMsg represents a tick for counter updates
type TickCounterMsg time.Time

// InitialCounterModel creates a new counter model
func InitialCounterModel() CounterModel {
	return CounterModel{
		count: 0,
	}
}

// Init implements the tea.Model interface - returns initial command
func (m CounterModel) Init() tea.Cmd {
	return doTick()
}

// Update implements the tea.Model interface - handles messages and updates state
func (m CounterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	
	case TickCounterMsg:
		m.count++
		return m, doTick()
	
	default:
		return m, nil
	}
}

// View implements the tea.Model interface - renders the current state
func (m CounterModel) View() string {
	if m.isQuitting {
		return fmt.Sprintf("Counter stopped at: %d\nThanks for watching!\n", m.count)
	}
	
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")). // Bright blue
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")) // Gray border
	
	content := fmt.Sprintf("Counting: %d\n\nPress 'q', 'esc', or 'ctrl+c' to quit", m.count)
	
	return style.Render(content)
}

// doTick returns a command that sends a TickCounterMsg after 1 second
func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickCounterMsg(t)
	})
}

// RunCounter starts a counter that increments every second
func RunCounter() error {
	model := InitialCounterModel()
	p := tea.NewProgram(model)
	
	_, err := p.Run()
	return err
}