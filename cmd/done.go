package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	iclient "github.com/marianoalvarez/freelow/internal/client"
	"github.com/marianoalvarez/freelow/internal/git"
	"github.com/marianoalvarez/freelow/internal/timer"
)

var doneCmd = &cobra.Command{
	Use:   "done [mensaje]",
	Short: "Para el timer, hace commit formateado y push",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noPush, _ := cmd.Flags().GetBool("no-push")

		sessions, err := timer.LoadSessions()
		if err != nil {
			return err
		}
		active := sessions.ActiveSession()
		if active == nil {
			return fmt.Errorf("no hay ninguna tarea activa")
		}

		now := time.Now()
		active.End = &now
		active.DurationMin = int(active.ElapsedDuration().Minutes())

		cfg, err := iclient.Load()
		if err != nil {
			return err
		}

		clientID := sessions.Client
		if client, err := cfg.GetActive(); err == nil {
			clientID = client.ID
		}

		msg := active.Task
		if len(args) > 0 && args[0] != "" {
			msg = args[0]
		}

		commitMsg := fmt.Sprintf("[%s] %s (%s)",
			clientID, msg, timer.FormatDuration(active.ElapsedDuration()),
		)

		var hash string
		if git.IsRepo() {
			hash, err = git.Commit(commitMsg)
			if err != nil {
				fmt.Printf("  (sin cambios para commitear)\n")
			} else {
				active.Commits = append(active.Commits, hash)
				if !noPush {
					if err := git.Push(); err != nil {
						fmt.Printf("  aviso: push falló: %v\n", err)
					}
				}
			}
		}

		if err := timer.SaveSessions(sessions); err != nil {
			return err
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf("✓ Tarea completada: %s", active.Task)))
		fmt.Printf("  Tiempo: %s\n", timer.FormatDuration(active.ElapsedDuration()))
		if hash != "" {
			fmt.Printf("  Commit: %s — %s\n", hash, commitMsg)
		}
		return nil
	},
}

func init() {
	doneCmd.Flags().Bool("no-push", false, "No hacer push automático")
}
