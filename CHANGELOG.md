# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-03-22

### Added
- `freelow client add/list/switch` — client management with persistence in `~/.freelow/clients.json`
- `freelow task "name"` — creates a `task/...` git branch and starts the timer
- `freelow pause` / `freelow resume` — pause and resume the active timer
- `freelow done [message]` — stops the timer, creates a formatted commit, and auto-pushes
- `freelow done --no-push` — same without push
- `freelow status` — shows active client, current task, and elapsed time
- `freelow report` — weekly markdown report ready to copy
- `freelow report --month` — monthly report
- `freelow init` — initializes `.freelow/sessions.json` in the current repo
- Auto-formatted commits: `[client] task (Xh Ymin)`
- Terminal colors with Lipgloss

[Unreleased]: https://github.com/marianoalvarez/freelow/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/marianoalvarez/freelow/releases/tag/v0.1.0
