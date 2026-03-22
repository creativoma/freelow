package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	iclient "github.com/creativoma/freelow/internal/client"
	"github.com/creativoma/freelow/internal/timer"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Muestra el cliente y tarea activos",
	RunE: func(cmd *cobra.Command, args []string) error {
		activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

		cfg, err := iclient.Load()
		if err != nil {
			return err
		}

		client, err := cfg.GetActive()
		if err != nil {
			fmt.Println(dimStyle.Render("Sin cliente activo. Usa: freelow client add <id>"))
		} else {
			fmt.Printf("%s %s\n", labelStyle.Render("Cliente:"), activeStyle.Render(client.Name))
		}

		sessions, err := timer.LoadSessions()
		if err != nil || sessions == nil {
			return nil
		}

		active := sessions.ActiveSession()
		if active == nil {
			fmt.Println(dimStyle.Render("Sin tarea activa."))
			return nil
		}

		elapsed := active.ElapsedDuration()
		status := "activa"
		if active.Paused {
			status = "pausada"
		}

		fmt.Printf("%s %s\n", labelStyle.Render("Tarea: "), activeStyle.Render(active.Task))
		fmt.Printf("%s %s\n", labelStyle.Render("Tiempo:"), timer.FormatDuration(elapsed))
		fmt.Printf("%s [%s]\n", labelStyle.Render("Estado:"), status)
		fmt.Printf("%s %s\n", labelStyle.Render("Rama:  "), active.Branch)
		return nil
	},
}
