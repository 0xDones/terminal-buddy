# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`clo` is a Go module (Go 1.25.6). The project is in early stages.

## Build & Development Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./path/to/package -run TestName

# Vet / lint
go vet ./...

# Format code
gofmt -w .

# Tidy dependencies
go mod tidy
```
