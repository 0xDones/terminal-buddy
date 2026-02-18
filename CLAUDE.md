# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`tb` (terminal-buddy) is a Go module (Go 1.25.6). The project is in early stages.

## Build & Development Commands

```bash
# Build (version auto-derived from git tags)
make build

# Build without make
go build -ldflags "-X main.version=$(git describe --tags --always --dirty)" -o tb .

# Run tests
go test ./...

# Run a single test
go test ./path/to/package -run TestName

# Vet / lint
make vet

# Format code
make fmt

# Tidy dependencies
make tidy
```

## CI/CD & Releases

Releases are fully automated via [Release Please](https://github.com/googleapis/release-please). No manual tagging required.

**How it works:** Merge PRs with [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, etc.) to main. Release Please opens a release PR that accumulates changes. Merging that PR auto-creates a git tag + GitHub Release, then GoReleaser builds and uploads cross-platform binaries.

```bash
# Validate GoReleaser config
goreleaser check

# Local dry-run (builds all platforms, doesn't publish)
goreleaser release --snapshot --clean
```

CI workflow (`.github/workflows/ci.yml`) runs `go vet`, `go test`, and `go build` on PRs and main pushes.
