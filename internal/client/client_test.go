package client

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToSlug(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"My Task", "my-task"},
		{"  fix login bug  ", "fix-login-bug"},
		{"Hello World!", "hello-world"},
		{"already-slug", "already-slug"},
		{"123 numbers", "123-numbers"},
		{"---trim---", "trim"},
	}
	for _, c := range cases {
		got := ToSlug(c.in)
		if got != c.want {
			t.Errorf("ToSlug(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestValidColor(t *testing.T) {
	valid := []string{"0", "4", "255", "#abc", "#ABC", "#aabbcc", "#AABBCC", "#000000", "#fff"}
	for _, c := range valid {
		if !ValidColor(c) {
			t.Errorf("ValidColor(%q) = false, want true", c)
		}
	}

	invalid := []string{"", "256", "-1", "blue", "red", "#gg0000", "#12345", "1.5", "#GGGGGG"}
	for _, c := range invalid {
		if ValidColor(c) {
			t.Errorf("ValidColor(%q) = true, want false", c)
		}
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".freelow", "clients.json")

	orig := os.Getenv("HOME")
	t.Setenv("HOME", tmp)
	defer t.Setenv("HOME", orig)

	// Load from non-existent file returns empty config.
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if len(cfg.Clients) != 0 {
		t.Errorf("expected empty clients, got %d", len(cfg.Clients))
	}

	// Save and reload.
	cfg.Clients = append(cfg.Clients, Client{ID: "acme", Name: "Acme Corp", Color: "4"})
	cfg.Active = "acme"
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load after save: %v", err)
	}
	if len(loaded.Clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(loaded.Clients))
	}
	if loaded.Clients[0].ID != "acme" {
		t.Errorf("ID = %q, want %q", loaded.Clients[0].ID, "acme")
	}
	if loaded.Active != "acme" {
		t.Errorf("Active = %q, want %q", loaded.Active, "acme")
	}
}

func TestGetActive(t *testing.T) {
	cfg := &Config{
		Active:  "b",
		Clients: []Client{{ID: "a"}, {ID: "b", Name: "B Corp"}},
	}

	c, err := cfg.GetActive()
	if err != nil {
		t.Fatalf("GetActive: %v", err)
	}
	if c.ID != "b" {
		t.Errorf("ID = %q, want %q", c.ID, "b")
	}
}

func TestGetActive_NoActive(t *testing.T) {
	cfg := &Config{}
	if _, err := cfg.GetActive(); err == nil {
		t.Error("expected error when no active client")
	}
}

func TestGetActive_NotFound(t *testing.T) {
	cfg := &Config{Active: "x", Clients: []Client{{ID: "a"}}}
	if _, err := cfg.GetActive(); err == nil {
		t.Error("expected error when active client not in list")
	}
}

func TestFindByID(t *testing.T) {
	cfg := &Config{Clients: []Client{{ID: "foo"}, {ID: "bar"}}}
	c, err := cfg.FindByID("bar")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if c.ID != "bar" {
		t.Errorf("ID = %q, want %q", c.ID, "bar")
	}
}

func TestFindByID_NotFound(t *testing.T) {
	cfg := &Config{Clients: []Client{{ID: "foo"}}}
	if _, err := cfg.FindByID("missing"); err == nil {
		t.Error("expected error for missing ID")
	}
}
