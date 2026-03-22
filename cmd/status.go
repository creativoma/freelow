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
	Short: "Show the active client and task",
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
			fmt.Println(dimStyle.Render("No active client. Run: freelow client add <id>"))
		} else {
			fmt.Printf("%s %s\n", labelStyle.Render("Client: "), activeStyle.Render(client.Name))
		}

		sessions, err := timer.LoadSessions()
		if err != nil || sessions == nil {
			return nil
		}

		active := sessions.ActiveSession()
		if active == nil {
			fmt.Println(dimStyle.Render("No active task."))
			return nil
		}

		elapsed := active.ElapsedDuration()
		status := "active"
		if active.Paused {
			status = "paused"
		}

		fmt.Printf("%s %s\n", labelStyle.Render("Task:  "), activeStyle.Render(active.Task))
		fmt.Printf("%s %s\n", labelStyle.Render("Time:  "), timer.FormatDuration(elapsed))
		fmt.Printf("%s [%s]\n", labelStyle.Render("Status:"), status)
		fmt.Printf("%s %s\n", labelStyle.Render("Branch:"), active.Branch)
		return nil
	},
}
