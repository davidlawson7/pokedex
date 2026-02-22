# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run ./cmd/main.go    # Run the application
go build ./cmd/         # Build the binary
go test ./...           # Run all tests
go test ./path/to/pkg   # Run tests for a specific package
go vet ./...            # Static analysis
```

## Architecture

This is an early-stage Go CLI application for managing Pokémon data.

**Module**: `github.com/davidlawson7/pokedex`

**Entry point**: `cmd/main.go`


