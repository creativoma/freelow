package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	iclient "github.com/creativoma/freelow/internal/client"
	"github.com/creativoma/freelow/internal/timer"
)

var initRepoCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize freelow in the current repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID, _ := cmd.Flags().GetString("client")

		if err := os.MkdirAll(".freelow", 0755); err != nil {
			return err
		}

		sessions, err := timer.LoadSessions()
		if err != nil {
			sessions = &timer.Sessions{}
		}

		if clientID != "" {
			sessions.Client = clientID
		} else {
			cfg, err := iclient.Load()
			if err == nil {
				if client, err := cfg.GetActive(); err == nil {
					sessions.Client = client.ID
				}
			}
		}

		if err := timer.SaveSessions(sessions); err != nil {
			return err
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(style.Render("✓ freelow initialized in this repo"))
		if sessions.Client != "" {
			fmt.Printf("  Client: %s\n", sessions.Client)
		}
		return nil
	},
}

func init() {
	initRepoCmd.Flags().StringP("client", "c", "", "Client ID for this repo")
}
