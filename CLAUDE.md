# call-my-agent

## What this is
A Go MCP server that turns Claude in Chrome into a personal job application agent.
Browse LinkedIn → Claude reads the page → calls this server → returns fit score,
cover letter, and saves the job to a local SQLite DB.

## Key design decisions
- Pure Go, no CGO (use modernc.org/sqlite, not mattn/go-sqlite3)
- MCP server runs over stdio (ServeStdio), not HTTP
- Scoring is prompt-based via Claude API, not hardcoded rules
- Candidate profile loaded from data/profile.json at startup
- No UI, no auth, no external dependencies beyond the Claude API

## What NOT to do
- Do not add HTTP endpoints — this is stdio MCP only
- Do not use CGO or cgo-dependent libraries
- Do not hardcode Rafael's profile in Go — load from profile.json
- Do not add a database until the tracking ticket (86c9wt0uw)

## Tool registry
| Tool                  | Status       | Ticket     |
|-----------------------|--------------|------------|
| get_candidate_profile | stub → real  | 86c9wt05a  |
| evaluate_job          | not started  | 86c9wt0gj  |
| generate_cover_letter | not started  | 86c9wt0nj  |
| save_job              | not started  | 86c9wt0uw  |
| list_jobs             | not started  | 86c9wt0uw  |

## Running locally
go run ./cmd/server
