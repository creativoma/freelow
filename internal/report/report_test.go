package report

import (
	"strings"
	"testing"
	"time"

	"github.com/creativoma/freelow/internal/timer"
)

var (
	base  = time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC)
	since = time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC)
	until = time.Date(2026, 3, 22, 23, 59, 59, 0, time.UTC)
)

func makeSession(task string, start time.Time, dur time.Duration, commits []string) timer.Session {
	end := start.Add(dur)
	return timer.Session{
		Task:    task,
		Branch:  "task/" + task,
		Start:   start,
		End:     &end,
		Commits: commits,
	}
}

func TestBuildFromSessions_Basic(t *testing.T) {
	sessions := []timer.Session{
		makeSession("fix-login", base, 80*time.Minute, []string{"abc", "def", "ghi"}),
		makeSession("header-redesign", base.Add(2*time.Hour), 125*time.Minute, []string{"jkl", "mno", "pqr", "stu", "vwx"}),
	}

	r := BuildFromSessions(sessions, "Acme Corp", since, until)

	if r.ClientName != "Acme Corp" {
		t.Errorf("ClientName = %q", r.ClientName)
	}
	if len(r.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(r.Tasks))
	}
	if r.Tasks[0].Name != "fix-login" {
		t.Errorf("Tasks[0].Name = %q", r.Tasks[0].Name)
	}
	if r.Tasks[0].Commits != 3 {
		t.Errorf("Tasks[0].Commits = %d, want 3", r.Tasks[0].Commits)
	}
	if r.Total != 205*time.Minute {
		t.Errorf("Total = %v, want 205min", r.Total)
	}
}

func TestBuildFromSessions_SkipsOpenSessions(t *testing.T) {
	sessions := []timer.Session{
		makeSession("done-task", base, 60*time.Minute, nil),
		{Task: "open-task", Start: base.Add(time.Hour)}, // no End
	}

	r := BuildFromSessions(sessions, "Acme", since, until)
	if len(r.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(r.Tasks))
	}
}

func TestBuildFromSessions_SkipsOutOfRange(t *testing.T) {
	old := time.Date(2026, 3, 1, 9, 0, 0, 0, time.UTC)
	sessions := []timer.Session{
		makeSession("old-task", old, 60*time.Minute, nil),
		makeSession("current-task", base, 30*time.Minute, nil),
	}

	r := BuildFromSessions(sessions, "Acme", since, until)
	if len(r.Tasks) != 1 || r.Tasks[0].Name != "current-task" {
		t.Errorf("expected only current-task, got %+v", r.Tasks)
	}
}

func TestBuildFromSessions_MergesSameTask(t *testing.T) {
	sessions := []timer.Session{
		makeSession("fix-login", base, 30*time.Minute, []string{"abc"}),
		makeSession("fix-login", base.Add(2*time.Hour), 30*time.Minute, []string{"def", "ghi"}),
	}

	r := BuildFromSessions(sessions, "Acme", since, until)
	if len(r.Tasks) != 1 {
		t.Fatalf("expected 1 merged task, got %d", len(r.Tasks))
	}
	if r.Tasks[0].Duration != 60*time.Minute {
		t.Errorf("merged duration = %v, want 60min", r.Tasks[0].Duration)
	}
	if r.Tasks[0].Commits != 3 {
		t.Errorf("merged commits = %d, want 3", r.Tasks[0].Commits)
	}
}

func TestGenerate_ContainsKeyFields(t *testing.T) {
	r := &Report{
		ClientName: "Acme Corp",
		Since:      since,
		Until:      until,
		Tasks: []TaskSummary{
			{Name: "fix-login", Duration: 80 * time.Minute, Commits: 3},
		},
		Total: 80 * time.Minute,
	}

	out, err := Generate(r, "weekly")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	checks := []string{"Acme Corp", "weekly", "fix-login", "1h 20min", "3"}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestGenerate_EmptyTasks(t *testing.T) {
	r := &Report{
		ClientName: "Acme Corp",
		Since:      since,
		Until:      until,
		Tasks:      nil,
		Total:      0,
	}

	_, err := Generate(r, "weekly")
	if err != nil {
		t.Fatalf("Generate with empty tasks returned error: %v", err)
	}
}
