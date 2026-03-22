package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestIsRepo_True(t *testing.T) {
	tmp := t.TempDir()
	cmd := exec.Command("git", "init", tmp)
	if err := cmd.Run(); err != nil {
		t.Skip("git not available")
	}

	orig, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	if !IsRepo() {
		t.Error("IsRepo() = false inside a git repo, want true")
	}
}

func TestIsRepo_False(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	if IsRepo() {
		t.Error("IsRepo() = true outside a git repo, want false")
	}
}

func TestRun_InvalidCommand(t *testing.T) {
	orig, _ := os.Getwd()
	tmp := t.TempDir()
	os.Chdir(tmp)
	defer os.Chdir(orig)

	_, err := Run("not-a-real-subcommand-xyz")
	if err == nil {
		t.Error("expected error for invalid git subcommand")
	}
}

func initRepo(t *testing.T) (dir string, cleanup func()) {
	t.Helper()
	tmp := t.TempDir()
	if err := exec.Command("git", "init", tmp).Run(); err != nil {
		t.Skip("git not available")
	}
	orig, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()
	return tmp, func() { os.Chdir(orig) }
}

func makeCommit(t *testing.T, tmp, msg string) {
	t.Helper()
	f := filepath.Join(tmp, "file.txt")
	os.WriteFile(f, []byte(msg), 0644)
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", msg).Run()
}

func TestCurrentBranch(t *testing.T) {
	tmp, cleanup := initRepo(t)
	defer cleanup()
	makeCommit(t, tmp, "init")

	branch, err := CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch: %v", err)
	}
	if branch == "" {
		t.Error("CurrentBranch returned empty string")
	}
}

func TestCreateBranch(t *testing.T) {
	tmp, cleanup := initRepo(t)
	defer cleanup()
	makeCommit(t, tmp, "init")

	if err := CreateBranch("task/my-feature"); err != nil {
		t.Fatalf("CreateBranch: %v", err)
	}

	branch, _ := CurrentBranch()
	if branch != "task/my-feature" {
		t.Errorf("CurrentBranch = %q, want %q", branch, "task/my-feature")
	}
}

func TestLog_ReturnsEntries(t *testing.T) {
	tmp, cleanup := initRepo(t)
	defer cleanup()
	makeCommit(t, tmp, "feat: add login")

	since := time.Now().Add(-1 * time.Hour)
	until := time.Now().Add(1 * time.Hour)

	entries, err := Log(since, until)
	if err != nil {
		t.Fatalf("Log: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected at least one log entry")
	}
	if entries[0].Hash == "" {
		t.Error("entry has empty hash")
	}
	if entries[0].Message == "" {
		t.Error("entry has empty message")
	}
}

func TestLog_OutOfRange(t *testing.T) {
	tmp, cleanup := initRepo(t)
	defer cleanup()
	makeCommit(t, tmp, "old commit")

	// Range in the future — should return no entries.
	since := time.Now().Add(1 * time.Hour)
	until := time.Now().Add(2 * time.Hour)

	entries, err := Log(since, until)
	if err != nil {
		t.Fatalf("Log: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
