package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "freelow",
	Short: "Tu copiloto freelance para git",
	Long:  "Tracking de horas, tareas y commits por cliente. Sin apps externas.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(initRepoCmd)
}
