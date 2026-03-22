package timer

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		input time.Duration
		want  string
	}{
		{0, "0min"},
		{45 * time.Minute, "45min"},
		{60 * time.Minute, "1h 00min"},
		{90 * time.Minute, "1h 30min"},
		{125 * time.Minute, "2h 05min"},
	}
	for _, c := range cases {
		got := FormatDuration(c.input)
		if got != c.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestElapsedDuration_Closed(t *testing.T) {
	start := time.Date(2026, 3, 22, 9, 0, 0, 0, time.UTC)
	end := start.Add(90 * time.Minute)
	sess := Session{Start: start, End: &end}

	got := sess.ElapsedDuration()
	if got != 90*time.Minute {
		t.Errorf("got %v, want 90min", got)
	}
}

func TestElapsedDuration_WithPause(t *testing.T) {
	start := time.Date(2026, 3, 22, 9, 0, 0, 0, time.UTC)
	end := start.Add(90 * time.Minute)
	sess := Session{
		Start:      start,
		End:        &end,
		PausedSecs: 10 * 60, // 10 min paused
	}

	got := sess.ElapsedDuration()
	if got != 80*time.Minute {
		t.Errorf("got %v, want 80min", got)
	}
}

func TestElapsedDuration_NeverNegative(t *testing.T) {
	start := time.Date(2026, 3, 22, 9, 0, 0, 0, time.UTC)
	end := start.Add(10 * time.Minute)
	sess := Session{
		Start:      start,
		End:        &end,
		PausedSecs: 60 * 60, // paused longer than elapsed
	}

	got := sess.ElapsedDuration()
	if got != 0 {
		t.Errorf("got %v, want 0", got)
	}
}

func TestActiveSession(t *testing.T) {
	end := time.Now()
	s := &Sessions{
		Sessions: []Session{
			{Task: "done", End: &end},
			{Task: "active"},
		},
	}

	active := s.ActiveSession()
	if active == nil {
		t.Fatal("expected active session, got nil")
	}
	if active.Task != "active" {
		t.Errorf("got task %q, want %q", active.Task, "active")
	}
}

func TestActiveSession_None(t *testing.T) {
	end := time.Now()
	s := &Sessions{
		Sessions: []Session{
			{Task: "done", End: &end},
		},
	}
	if got := s.ActiveSession(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestActiveSession_Empty(t *testing.T) {
	s := &Sessions{}
	if got := s.ActiveSession(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}
