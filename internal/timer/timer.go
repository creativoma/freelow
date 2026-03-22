package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Session struct {
	Task        string     `json:"task"`
	Branch      string     `json:"branch"`
	Start       time.Time  `json:"start"`
	End         *time.Time `json:"end,omitempty"`
	DurationMin int        `json:"duration_min,omitempty"`
	Commits     []string   `json:"commits,omitempty"`
	Paused      bool       `json:"paused,omitempty"`
	PausedAt    *time.Time `json:"paused_at,omitempty"`
	PausedSecs  int        `json:"paused_secs,omitempty"`
}

type Sessions struct {
	Client   string    `json:"client"`
	Sessions []Session `json:"sessions"`
}

func SessionsPath() string {
	return filepath.Join(".freelow", "sessions.json")
}

func LoadSessions() (*Sessions, error) {
	data, err := os.ReadFile(SessionsPath())
	if os.IsNotExist(err) {
		return &Sessions{}, nil
	}
	if err != nil {
		return nil, err
	}
	var sessions Sessions
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, err
	}
	return &sessions, nil
}

func SaveSessions(s *Sessions) error {
	if err := os.MkdirAll(".freelow", 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(SessionsPath(), data, 0644)
}

func (s *Sessions) ActiveSession() *Session {
	for i := range s.Sessions {
		if s.Sessions[i].End == nil {
			return &s.Sessions[i]
		}
	}
	return nil
}

func (sess *Session) ElapsedDuration() time.Duration {
	end := time.Now()
	if sess.End != nil {
		end = *sess.End
	}
	elapsed := end.Sub(sess.Start)
	elapsed -= time.Duration(sess.PausedSecs) * time.Second
	if sess.Paused && sess.PausedAt != nil {
		elapsed -= time.Since(*sess.PausedAt)
	}
	if elapsed < 0 {
		return 0
	}
	return elapsed
}

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dmin", h, m)
	}
	return fmt.Sprintf("%dmin", m)
}
