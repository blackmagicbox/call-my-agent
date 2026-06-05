# Contributing / Developer Guide

This document covers everything you need to work on the codebase — architecture decisions, project layout, how to run tests, and the contribution workflow.

---

## Architecture overview

Call My Agent is a **stdio MCP server** written in pure Go. There is no HTTP layer, no daemon, and no UI. Claude connects to the process over stdin/stdout using the [Model Context Protocol](https://modelcontextprotocol.io/).

Key design decisions:

- **Pure Go, no CGO** — uses `modernc.org/sqlite` (not `mattn/go-sqlite3`) so the binary cross-compiles without a C toolchain.
- **LLM provider abstracted** — all LLM calls go through `llm.Provider`. Swap providers by changing `LLM_PROVIDER`. Adding a new provider means implementing one interface.
- **Prompts as files** — prompt strings live in `internal/prompts/` as embedded `.txt`/`.tmpl` files, not as Go string literals. This keeps prompts reviewable and diffable.
- **Profile loaded at startup** — `data/profile.json` is read once when the server starts. No hot-reload.

---

## Project layout

```
cmd/server/
└── main.go                      ← entrypoint: wires DB, profile, LLM, tools, MCP server

internal/
├── llm/
│   ├── llm.go                   ← Provider interface
│   ├── gemini.go                ← GeminiProvider implementation
│   ├── config.go                ← FromConfig() reads env vars
│   ├── mock.go                  ← MockProvider for use in tool tests
│   └── llm_test.go              ← config + mock tests
├── prompts/
│   ├── prompts.go               ← go:embed directives + Render* functions
│   ├── evaluate_job.system.txt
│   ├── evaluate_job.user.tmpl
│   ├── cover_letter.system.tmpl
│   ├── cover_letter.user.tmpl
│   └── prompts_test.go
├── profile/
│   └── profile.go               ← CandidateProfile struct
├── job/
│   └── job.go                   ← Job + Compensation structs
├── tools/
│   ├── profile.go               ← get_candidate_profile handler
│   ├── evaluate.go              ← evaluate_job tool + handler
│   ├── coverletter.go           ← generate_cover_letter tool + handler
│   ├── tracking.go              ← save_job / list_jobs / get_application_status
│   ├── utils.go                 ← stripJSONFences() shared helper
│   └── *_test.go                ← one test file per handler
├── db/
│   ├── db.go                    ← SQLite layer (Open, SaveJob, ListJobs, GetStatusSummary)
│   └── db_test.go
└── testutil/
    └── testutil.go              ← LoadFixture + MustUnmarshal[T] for tests

testdata/
├── profile.json                 ← stub profile for tests (safe to commit — fake data)
└── job.json                     ← stub job for tests

data/
├── profile.json                 ← your real profile (gitignored)
├── profile.example.json         ← anonymised template (committed)
└── jobs.db                      ← SQLite database (gitignored, created at runtime)
```

---

## Environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LLM_PROVIDER` | Yes | — | `gemini` (only supported provider) |
| `LLM_API_KEY` | Yes | — | API key for the provider |
| `LLM_MODEL` | No | `gemini-2.5-flash` | Model name |

---

## Running locally

```bash
# Run the server directly (reads data/profile.json and data/jobs.db)
LLM_PROVIDER=gemini LLM_API_KEY=your-key go run ./cmd/server

# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run a single package
go test ./internal/tools/...
```

---

## Adding a new tool

1. **Define the struct/schema** in `internal/job/` or `internal/profile/` if needed.
2. **Write the prompt** — add `.txt` (static) or `.tmpl` (dynamic) files to `internal/prompts/` and wire them up in `prompts.go`.
3. **Implement the handler** — create `internal/tools/yourtool.go` with a `YourTool() mcp.Tool` function and a `HandleYourTool(...) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)` handler. Never call the LLM directly — always use `llm.Provider`.
4. **Write tests** — use `MockProvider` from `internal/llm/mock.go` and fixtures from `testutil.LoadFixture`. Aim for happy path + at least two error cases.
5. **Register in main.go** — add `s.AddTool(YourTool(), HandleYourTool(...))`.

---

## Adding a new LLM provider

Implement the `llm.Provider` interface:

```go
type Provider interface {
    Generate(ctx context.Context, system, user string) (string, error)
}
```

Then add a case for your provider in `llm.FromConfig()`.

---

## Test conventions

- Tool tests use `MockProvider` — never make real LLM calls in tests.
- DB tests use `db.Open(":memory:")` — never touch the filesystem.
- Shared fixture loading goes through `testutil.LoadFixture(t, "profile.json")`.
- `testutil.MustUnmarshal[T](t, data)` for unmarshalling fixtures into typed structs.
- All error returns must be handled — no bare discards. Use `_ =` only with an explanatory comment.
- Test code is held to the same standard as production code.

---

## Contribution workflow

This repo enforces PRs on `main`. Direct pushes are blocked by both GitHub rulesets and a local pre-push hook (installed automatically if you run setup).

```bash
# Start a new piece of work
git checkout -b feat/your-feature-name

# ... make changes, write tests ...
go test ./...

# Push and open a PR
git push origin feat/your-feature-name
gh pr create
```

Branch naming conventions:

| Prefix | Use for |
|--------|---------|
| `feat/` | New features |
| `fix/` | Bug fixes |
| `chore/` | Tooling, deps, config |
| `docs/` | Documentation only |
| `test/` | Tests only |

CI runs `go test ./...` and `golangci-lint` on every push and PR. Both must pass before merging.

---

## Code quality baseline

This project targets Google SRE L5/L6 as a portfolio artifact. That means:

- `go vet`, `staticcheck`, and `errcheck` findings are treated as bugs, not style nits.
- Lint violations block CI.
- Prefer explicit over implicit — no magic, no hidden behaviour.
- Prompts are files, not strings. Configs are env vars, not hardcoded values.
