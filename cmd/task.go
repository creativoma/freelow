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

var taskCmd = &cobra.Command{
	Use:   "task [nombre]",
	Short: "Crea una tarea (rama git + timer) o lista tareas",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag || len(args) == 0 {
			return listTasks()
		}
		return createTask(args[0])
	},
}

func createTask(name string) error {
	cfg, err := iclient.Load()
	if err != nil {
		return err
	}
	client, err := cfg.GetActive()
	if err != nil {
		return fmt.Errorf("no hay cliente activo. Usa: freelow client add <id>")
	}

	sessions, err := timer.LoadSessions()
	if err != nil {
		return err
	}
	if active := sessions.ActiveSession(); active != nil {
		return fmt.Errorf("ya hay una tarea activa: %s. Usa 'freelow done' o 'freelow pause'", active.Task)
	}

	slug := iclient.ToSlug(name)
	branch := fmt.Sprintf("task/%s", slug)

	if git.IsRepo() {
		if err := git.CreateBranch(branch); err != nil {
			fmt.Printf("  (aviso: no se pudo crear rama %s: %v)\n", branch, err)
		}
	}

	now := time.Now()
	sessions.Client = client.ID
	sessions.Sessions = append(sessions.Sessions, timer.Session{
		Task:   slug,
		Branch: branch,
		Start:  now,
	})

	if err := timer.SaveSessions(sessions); err != nil {
		return err
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	fmt.Println(style.Render(fmt.Sprintf("▶ Tarea iniciada: %s [%s]", name, client.ID)))
	fmt.Printf("  Rama:   %s\n", branch)
	fmt.Printf("  Inicio: %s\n", now.Format("15:04:05"))
	return nil
}

func listTasks() error {
	sessions, err := timer.LoadSessions()
	if err != nil {
		return err
	}

	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	var open []*timer.Session
	for i := range sessions.Sessions {
		if sessions.Sessions[i].End == nil {
			open = append(open, &sessions.Sessions[i])
		}
	}

	if len(open) == 0 {
		fmt.Println(dimStyle.Render("No hay tareas abiertas."))
		return nil
	}

	fmt.Println("Tareas abiertas:")
	for _, s := range open {
		elapsed := timer.FormatDuration(s.ElapsedDuration())
		status := "activa"
		if s.Paused {
			status = "pausada"
		}
		line := fmt.Sprintf("  %-30s  %s  [%s]", s.Task, elapsed, status)
		if !s.Paused {
			fmt.Println(activeStyle.Render(line))
		} else {
			fmt.Println(dimStyle.Render(line))
		}
	}
	return nil
}

func init() {
	taskCmd.Flags().BoolP("list", "l", false, "Lista tareas abiertas")
}
