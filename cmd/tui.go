/*
Copyright © 2025 dreamsofcode

*/
package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/dreamsofcode-io/cli-cms/internal/tui"
	"github.com/dreamsofcode-io/cli-cms/internal/ui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Demonstrates TUI components using Bubble Tea",
	Long: `This command demonstrates Terminal User Interface (TUI) components
using the Bubble Tea framework. It includes examples of spinners and progress bars.`,
}

// spinnerCmd demonstrates a spinner component
var spinnerCmd = &cobra.Command{
	Use:   "spinner",
	Short: "Show a spinning indicator",
	Long:  `Displays a spinning indicator with a message. Press 'q', 'esc', or 'ctrl+c' to quit.`,
	RunE:  runSpinner,
}

// progressCmd demonstrates a progress bar component
var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show a progress bar",
	Long:  `Displays an animated progress bar that automatically completes. Press 'q', 'esc', or 'ctrl+c' to cancel.`,
	RunE:  runProgress,
}

// counterCmd demonstrates a counter with ticking
var counterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Show a counter that increments every second",
	Long:  `Displays a counter that increments every second, demonstrating async updates. Press 'q', 'esc', or 'ctrl+c' to quit.`,
	RunE:  runCounter,
}

// taskCmd demonstrates running tasks with spinners
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Run a simulated task with spinner",
	Long:  `Runs a simulated long-running task with a spinner. Use --fail to simulate a failing task.`,
	RunE:  runTask,
}

var (
	taskDuration time.Duration
	taskFail     bool
)

func runSpinner(cmd *cobra.Command, args []string) error {
	message := "Loading posts from database..."
	if len(args) > 0 {
		message = args[0]
	}

	ui.PrintInfo("Starting spinner demo...\n")
	fmt.Printf("Use 'q', 'esc', or 'ctrl+c' to quit\n\n")

	// Try to run full TUI, fallback to demo if TTY not available
	err := tui.RunSpinner(message)
	if err != nil {
		ui.PrintWarning("TTY not available, showing demo instead...\n")
		tui.DemoSpinner(message, 3*time.Second)
		return nil
	}

	ui.PrintSuccess("Spinner demo completed!\n")
	return nil
}

func runProgress(cmd *cobra.Command, args []string) error {
	message := "Processing blog posts..."
	if len(args) > 0 {
		message = args[0]
	}

	ui.PrintInfo("Starting progress bar demo...\n")
	fmt.Printf("The progress bar will auto-complete, or use 'q', 'esc', or 'ctrl+c' to cancel\n\n")

	// Try to run full TUI, fallback to demo if TTY not available
	err := tui.RunProgress(message)
	if err != nil {
		ui.PrintWarning("TTY not available, showing demo instead...\n")
		tui.DemoProgress(message)
		return nil
	}

	ui.PrintSuccess("Progress bar demo completed!\n")
	return nil
}

func runCounter(cmd *cobra.Command, args []string) error {
	ui.PrintInfo("Starting counter demo...\n")
	fmt.Printf("The counter will increment every second. Use 'q', 'esc', or 'ctrl+c' to quit\n\n")

	// Try to run full TUI, fallback to demo if TTY not available
	err := tui.RunCounter()
	if err != nil {
		ui.PrintWarning("TTY not available, showing demo instead...\n")
		tui.DemoCounter(5 * time.Second)
		return nil
	}

	ui.PrintSuccess("Counter demo completed!\n")
	return nil
}

func runTask(cmd *cobra.Command, args []string) error {
	message := "Processing data"
	if len(args) > 0 {
		message = args[0]
	}

	ui.PrintInfo("Starting task demo...\n")
	fmt.Printf("Running task: %s (duration: %v)\n\n", message, taskDuration)

	// Create the task
	var task func() error
	if taskFail {
		task = tui.SimulateFailingTask(message, taskDuration, errors.New("simulated task failure"))
	} else {
		task = tui.SimulateTask(message, taskDuration)
	}

	// Try to run full TUI, fallback to simple execution if TTY not available
	err := tui.RunTask(message, task)
	if err != nil {
		// Check if it's a TTY error or actual task error
		if err.Error() == "simulated task failure" {
			ui.PrintError("Task failed: %v\n", err)
			return nil
		}
		
		// TTY not available, run task directly
		ui.PrintWarning("TTY not available, running task directly...\n")
		ui.PrintInfo("Executing: %s...\n", message)
		
		err = task()
		if err != nil {
			ui.PrintError("Task failed: %v\n", err)
		} else {
			ui.PrintSuccess("Task completed successfully!\n")
		}
		return nil
	}

	ui.PrintSuccess("Task demo completed!\n")
	return nil
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.AddCommand(spinnerCmd)
	tuiCmd.AddCommand(progressCmd)
	tuiCmd.AddCommand(counterCmd)
	tuiCmd.AddCommand(taskCmd)

	// Task command flags
	taskCmd.Flags().DurationVarP(&taskDuration, "duration", "t", 3*time.Second, "Duration of the simulated task")
	taskCmd.Flags().BoolVar(&taskFail, "fail", false, "Make the task fail after the duration")
}