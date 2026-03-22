package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/creativoma/freelow/internal/timer"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the active timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions, err := timer.LoadSessions()
		if err != nil {
			return err
		}
		active := sessions.ActiveSession()
		if active == nil {
			return fmt.Errorf("no active task")
		}
		if active.Paused {
			return fmt.Errorf("task is already paused")
		}
		now := time.Now()
		active.Paused = true
		active.PausedAt = &now
		if err := timer.SaveSessions(sessions); err != nil {
			return err
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf("⏸  Paused: %s", active.Task)))
		fmt.Printf("  Elapsed: %s\n", timer.FormatDuration(active.ElapsedDuration()))
		return nil
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the paused timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions, err := timer.LoadSessions()
		if err != nil {
			return err
		}
		active := sessions.ActiveSession()
		if active == nil {
			return fmt.Errorf("no active task")
		}
		if !active.Paused {
			return fmt.Errorf("task is not paused")
		}
		if active.PausedAt != nil {
			active.PausedSecs += int(time.Since(*active.PausedAt).Seconds())
		}
		active.Paused = false
		active.PausedAt = nil
		if err := timer.SaveSessions(sessions); err != nil {
			return err
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf("▶ Resumed: %s", active.Task)))
		return nil
	},
}
