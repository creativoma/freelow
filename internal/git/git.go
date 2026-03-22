package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

func CurrentBranch() (string, error) {
	return Run("rev-parse", "--abbrev-ref", "HEAD")
}

func CreateBranch(name string) error {
	_, err := Run("checkout", "-b", name)
	return err
}

func Commit(message string) (string, error) {
	if _, err := Run("add", "-A"); err != nil {
		return "", err
	}
	if _, err := Run("commit", "-m", message); err != nil {
		return "", err
	}
	return Run("rev-parse", "--short", "HEAD")
}

func Push() error {
	branch, err := CurrentBranch()
	if err != nil {
		return err
	}
	_, err = Run("push", "-u", "origin", branch)
	return err
}

type LogEntry struct {
	Hash    string
	Message string
	Date    time.Time
}

func Log(since, until time.Time) ([]LogEntry, error) {
	args := []string{
		"log",
		"--pretty=format:%H\t%s\t%ai",
		"--after=" + since.Format(time.RFC3339),
		"--before=" + until.Format(time.RFC3339),
	}
	out, err := Run(args...)
	if err != nil || out == "" {
		return nil, err
	}
	var entries []LogEntry
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}
		t, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[2])
		entries = append(entries, LogEntry{
			Hash:    parts[0],
			Message: parts[1],
			Date:    t,
		})
	}
	return entries, nil
}

func IsRepo() bool {
	_, err := Run("rev-parse", "--git-dir")
	return err == nil
}
