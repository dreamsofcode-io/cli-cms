package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TaskModel represents a long-running task with a spinner
type TaskModel struct {
	spinner    spinner.Model
	isQuitting bool
	message    string
	task       func() error
	result     error
	completed  bool
}

// TaskCompleteMsg represents task completion
type TaskCompleteMsg struct {
	err error
}

// InitialTaskModel creates a new task model
func InitialTaskModel(message string, task func() error) TaskModel {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
	
	return TaskModel{
		spinner: s,
		message: message,
		task:    task,
	}
}

// Init implements the tea.Model interface
func (m TaskModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runTask())
}

// Update implements the tea.Model interface
func (m TaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	
	case TaskCompleteMsg:
		m.completed = true
		m.result = msg.err
		m.isQuitting = true
		return m, tea.Quit
	
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// View implements the tea.Model interface
func (m TaskModel) View() string {
	if m.isQuitting {
		if m.completed {
			if m.result != nil {
				return fmt.Sprintf("❌ %s failed: %v\n", m.message, m.result)
			}
			return fmt.Sprintf("✅ %s completed successfully!\n", m.message)
		}
		return fmt.Sprintf("❌ %s cancelled\n", m.message)
	}
	
	return fmt.Sprintf("%s %s... (Press q to cancel)", m.spinner.View(), m.message)
}

// runTask returns a command that executes the task
func (m TaskModel) runTask() tea.Cmd {
	return func() tea.Msg {
		err := m.task()
		return TaskCompleteMsg{err: err}
	}
}

// RunTask executes a long-running task with a spinner
func RunTask(message string, task func() error) error {
	model := InitialTaskModel(message, task)
	p := tea.NewProgram(model)
	
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	
	// Return the task result
	if taskModel, ok := finalModel.(TaskModel); ok {
		return taskModel.result
	}
	
	return nil
}

// SimulateTask creates a task that simulates work with a delay
func SimulateTask(name string, duration time.Duration) func() error {
	return func() error {
		time.Sleep(duration)
		return nil
	}
}

// SimulateFailingTask creates a task that fails after a delay
func SimulateFailingTask(name string, duration time.Duration, err error) func() error {
	return func() error {
		time.Sleep(duration)
		return err
	}
}