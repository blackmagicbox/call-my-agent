# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Call My Agent is an MCP (Model Context Protocol) server written in Go. It serves candidate profile data to Claude (via the Chrome extension or Claude Desktop) and is designed to help users streamline their job search — evaluating job fitness against user-defined criteria and generating tailored cover letters.

## Build & Run

```bash
go build -o call-my-agent ./cmd/server
go run ./cmd/server
```

The server communicates over **stdio** (stdin/stdout) using the MCP protocol. It is meant to be launched by a host application (e.g. Claude Desktop), not run as a standalone HTTP server.

## Test

```bash
go test ./...                  # all tests
go test ./internal/...         # internal packages only
go test -run TestName ./path   # single test
```

## Architecture

- **`cmd/server/main.go`** — Entrypoint. Loads `data/profile.json`, registers MCP tools, and starts the stdio server.
- **`internal/profile/`** — Data structures for candidate profiles (experience, education, languages, preferences).
- **`internal/job/`** — Data structures for job postings and compensation.
- **`internal/tools/`** — Stubs for future MCP tool handlers (evaluate, coverletter, tracking).
- **`internal/db/`** — Stub for future persistence layer.
- **`data/`** — JSON files: `profile.json` (user profile) and `profile.example.json` (template).

The project follows the standard Go layout: `cmd/` for binaries, `internal/` for implementation details.

## MCP Tools

| Tool | Status | Description |
|---|---|---|
| `get_candidate_profile` | Implemented | Returns the candidate's full profile |
| `evaluate_job` | Planned | Evaluate a job listing against user preferences |
| `generate_cover_letter` | Planned | Generate a tailored cover letter |
| `track_application` | Planned | Track job application status |

## Go Module

Module path: `github.com/blackmagicbox/call-my-agent`

Key dependency: [`mcp-go`](https://github.com/mark3labs/mcp-go) for the MCP protocol implementation.
