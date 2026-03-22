# freelow

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)


CLI for freelancers with multiple clients. Track hours, tasks, and commits from the terminal — no external apps required.

Lives inside git and records everything automatically as you work.

---

## Installation

### From source

```bash
git clone https://github.com/creativoma/freelow.git
cd freelow
go install .
```

### With `go install`

```bash
go install github.com/creativoma/freelow@latest
```

**Requirements:** Go 1.21+

---

## Quick start

```bash
# 1. Add a client
freelow client add acme-corp

# 2. Initialize freelow in your repo
cd ~/projects/acme-project
freelow init

# 3. Start working
freelow task "fix login bug"

# 4. Check how long you've been going
freelow status

# 5. Finish (auto commit + push)
freelow done "fixed Google auth error"

# 6. Weekly report ready to paste
freelow report
```

---

## Commands

### Clients

```bash
freelow client add <id>                    # register a new client
freelow client add <id> -n "Name"          # with a custom display name
freelow client add <id> -c 4               # with an ANSI color (0-255)
freelow client add <id> -c "#e879f9"       # with a hex color (#RGB or #RRGGBB)
freelow client list                        # list all clients
freelow client switch <id>                 # set the active client
```

### Tasks

```bash
freelow task "description"          # create task/... branch + start timer
freelow task                        # list open tasks
freelow task -l                     # same
freelow pause                       # pause the timer
freelow resume                      # resume the timer
```

### Closing work

```bash
freelow done                        # stop timer, commit, push
freelow done "commit message"       # with a custom message
freelow done --no-push              # without auto push
```

Commit message is formatted automatically:
```
[acme-corp] fix-login-bug (1h 20min)
```

### Status & reports

```bash
freelow status                      # active client, current task, elapsed time
freelow report                      # weekly report for the active client
freelow report --month              # monthly report
freelow report acme-corp --week     # report for a specific client
```

### Repo setup

```bash
freelow init                        # initialize .freelow/ in the current repo
freelow init -c acme-corp           # link to a specific client
```

---

## Sample report output

```markdown
## Weekly report · Acme Corp · 17–22 Mar 2026

**Total hours:** 6h 40min

| Task              | Time      | Commits |
|-------------------|-----------|---------|
| fix-login-bug     | 1h 20min  | 3       |
| header-redesign   | 2h 05min  | 5       |
| fix-cart-bug      | 3h 15min  | 4       |

### Changes
- Fixed Google authentication error on login
- New responsive header for mobile
- Fixed total calculation bug in cart
```

Markdown ready to paste into an email or Notion.

---

## Data model

```
~/.freelow/clients.json     — global client list
.freelow/sessions.json      — per-repo sessions (already in .gitignore)
```

---

## Stack

| Library | Purpose |
|---|---|
| [Cobra](https://github.com/spf13/cobra) | Subcommands and flags |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal colors |
| `os/exec` | Run git commands |
| `encoding/json` | Data persistence |
| `text/template` | Markdown report generation |

---

## Testing

```bash
go test ./internal/...          # run all tests
go test ./internal/... -cover   # with coverage
go test ./internal/... -v       # verbose
```

Tests cover `internal/timer` (duration formatting, elapsed calculation, pause logic), `internal/report` (session parsing, report generation), `internal/client` (config I/O, slug, color validation), and `internal/git` (repo detection, branch creation, log parsing).

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Author

**Mariano Alvarez** — [@creativoma](https://github.com/creativoma)

---

## License

MIT — see [LICENSE](LICENSE).
