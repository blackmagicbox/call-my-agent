# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Call My Agent is an MCP (Model Context Protocol) server written in Go. It integrates with the Claude chrome extension to assist users with job searching — evaluating job fitness against user-defined criteria and generating tailored cover letters.

## Build & Run

```bash
go build -o call-my-agent ./cmd/server
go run ./cmd/server
```

## Test

```bash
go test ./...                  # all tests
go test ./internal/...         # internal packages only
go test -run TestName ./path   # single test
```

## Architecture

- **`cmd/server/`** — Application entrypoint (`main.go`). Starts the MCP server.
- **`internal/`** — Private packages (not importable by external modules). All core logic lives here.

The project follows the standard Go layout: `cmd/` for binaries, `internal/` for implementation details.

## Go Module

Module path: `github.com/blackmagicbox/call-my-agent`
