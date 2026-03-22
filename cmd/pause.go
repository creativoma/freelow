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
	Short: "Pausa el timer activo",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions, err := timer.LoadSessions()
		if err != nil {
			return err
		}
		active := sessions.ActiveSession()
		if active == nil {
			return fmt.Errorf("no hay ninguna tarea activa")
		}
		if active.Paused {
			return fmt.Errorf("la tarea ya está pausada")
		}
		now := time.Now()
		active.Paused = true
		active.PausedAt = &now
		if err := timer.SaveSessions(sessions); err != nil {
			return err
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf("⏸  Pausada: %s", active.Task)))
		fmt.Printf("  Tiempo acumulado: %s\n", timer.FormatDuration(active.ElapsedDuration()))
		return nil
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Reanuda el timer pausado",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions, err := timer.LoadSessions()
		if err != nil {
			return err
		}
		active := sessions.ActiveSession()
		if active == nil {
			return fmt.Errorf("no hay ninguna tarea activa")
		}
		if !active.Paused {
			return fmt.Errorf("la tarea no está pausada")
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
		fmt.Println(style.Render(fmt.Sprintf("▶ Reanudada: %s", active.Task)))
		return nil
	},
}
