package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	iclient "github.com/creativoma/freelow/internal/client"
	"github.com/creativoma/freelow/internal/report"
	"github.com/creativoma/freelow/internal/timer"
)

var reportCmd = &cobra.Command{
	Use:   "report [client]",
	Short: "Generate a hours and tasks report",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		monthFlag, _ := cmd.Flags().GetBool("month")

		cfg, err := iclient.Load()
		if err != nil {
			return err
		}

		now := time.Now()
		var since, until time.Time
		var periodLabel string

		if monthFlag {
			since = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			until = now
			periodLabel = "monthly"
		} else {
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			since = now.AddDate(0, 0, -(weekday - 1))
			since = time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, since.Location())
			until = now
			periodLabel = "weekly"
		}

		var clientID string
		if len(args) > 0 {
			clientID = args[0]
		} else {
			client, err := cfg.GetActive()
			if err != nil {
				return err
			}
			clientID = client.ID
		}

		client, err := cfg.FindByID(clientID)
		if err != nil {
			return err
		}

		sessions, err := timer.LoadSessions()
		if err != nil {
			return err
		}

		r := report.BuildFromSessions(sessions.Sessions, client.Name, since, until)
		out, err := report.Generate(r, periodLabel)
		if err != nil {
			return err
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	reportCmd.Flags().Bool("week", false, "Weekly report (default)")
	reportCmd.Flags().Bool("month", false, "Monthly report")
	reportCmd.Flags().Bool("all", false, "All clients")
}
