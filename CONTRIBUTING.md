# Contributing to freelow

Thanks for taking the time to contribute.

## Getting started

```bash
git clone https://github.com/creativoma/freelow.git
cd freelow
go build ./...
go test ./...
```

## How to contribute

1. **Open an issue first** for non-trivial changes — describe what you want to fix or add before writing code
2. Fork the repo and create a branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add or update tests if relevant
5. Make sure everything passes: `go test ./... && go vet ./...`
6. Open a Pull Request against `main`

## Guidelines

- Keep PRs focused — one thing per PR
- Match the existing code style (`gofmt -w .` before committing)
- Commit messages in lowercase imperative: `add report --copy flag`, `fix timer negative duration`
- Don't add dependencies without discussing it first

## Running tests

```bash
go test ./internal/...          # unit tests
go test ./internal/... -cover   # with coverage
go vet ./...                    # static analysis
```

## Reporting bugs

Open an issue at [github.com/creativoma/freelow/issues](https://github.com/creativoma/freelow/issues) with:

- What you did
- What you expected
- What actually happened
- Your OS and Go version (`go version`)

## License

By contributing you agree that your changes will be licensed under the [MIT License](LICENSE).
