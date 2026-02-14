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
