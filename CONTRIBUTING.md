# Contributing to Cencori Go SDK

## Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/DanielPopoola/cencori-go.git
cd cencori-go
```

2. **Install development tools**
```bash
make install-tools
```

3. **Run tests**
```bash
make test
```

## Before Submitting PR

Run the full verification suite:
```bash
make verify
```

This runs:
- Code formatting (`gofmt`, `goimports`)
- Linting (`golangci-lint`)
- Tests with race detection
- Build verification

## Code Standards

### Formatting
- Use `gofmt -s` for simplification
- Run `goimports` to organize imports
- Max line length: 120 characters

### Linting
We use `golangci-lint` with strict rules:
```bash
make lint
```

Fix issues before committing.

### Testing
- **Unit tests**: All new code must have tests
- **Race detection**: Always run with `-race`
- **Coverage**: Maintain >80% coverage

```bash
make test-coverage
```

### Commit Messages
Follow conventional commits:
```
feat: add streaming support for chat
fix: handle context cancellation in streaming
docs: update README with examples
test: add integration tests for projects API
```

## Adding Examples

1. Create new directory: `examples/XX-feature-name/`
2. Add `main.go` with complete example
3. Update `examples/README.md`
4. Test it runs: `go run examples/XX-feature-name/main.go`

## Running CI Locally

Simulate GitHub Actions:
```bash
# Lint check
make lint

# Tests with coverage
make test

# Build
make build

# All checks
make verify
```

## Questions?

- Open an issue for bugs
- Start a discussion for features
- Join our Discord for quick questions