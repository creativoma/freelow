package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	iclient "github.com/creativoma/freelow/internal/client"
	"github.com/creativoma/freelow/internal/git"
	"github.com/creativoma/freelow/internal/timer"
)

var taskCmd = &cobra.Command{
	Use:   "task [name]",
	Short: "Start a task (git branch + timer) or list tasks",
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
		return fmt.Errorf("no active client. Run: freelow client add <id>")
	}

	sessions, err := timer.LoadSessions()
	if err != nil {
		return err
	}
	if active := sessions.ActiveSession(); active != nil {
		return fmt.Errorf("task already active: %s. Run 'freelow done' or 'freelow pause'", active.Task)
	}

	slug := iclient.ToSlug(name)
	branch := fmt.Sprintf("task/%s", slug)

	if git.IsRepo() {
		if err := git.CreateBranch(branch); err != nil {
			fmt.Printf("  (warning: could not create branch %s: %v)\n", branch, err)
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
	fmt.Println(style.Render(fmt.Sprintf("▶ Task started: %s [%s]", name, client.ID)))
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
		fmt.Println(dimStyle.Render("No open tasks."))
		return nil
	}

	fmt.Println("Open tasks:")
	for _, s := range open {
		elapsed := timer.FormatDuration(s.ElapsedDuration())
		status := "active"
		if s.Paused {
			status = "paused"
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
	taskCmd.Flags().BoolP("list", "l", false, "List open tasks")
}
