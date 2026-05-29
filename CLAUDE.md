# call-my-agent

## What this is
A Go MCP server that turns Claude in Chrome into a personal job application agent.
Browse LinkedIn → Claude reads the page → calls this server → returns fit score,
cover letter, and saves the job to a local SQLite DB.

## Key design decisions
- Pure Go, no CGO (use modernc.org/sqlite, not mattn/go-sqlite3)
- MCP server runs over stdio (ServeStdio), not HTTP
- Scoring is prompt-based via LLM, not hardcoded rules
- LLM provider abstracted behind llm.Provider interface — swap by changing LLM_PROVIDER env var
- Prompts live in internal/prompts/ as embedded files (.txt static, .tmpl dynamic)
- Candidate profile loaded from data/profile.json at startup (gitignored)
- No UI, no auth, no external dependencies beyond the LLM API

## What NOT to do
- Do not add HTTP endpoints — this is stdio MCP only
- Do not use CGO or cgo-dependent libraries
- Do not hardcode the candidate profile in Go — load from data/profile.json
- Do not add a database until the tracking ticket (86c9wt0uw)
- Do not put inline prompt strings in tool code — use internal/prompts package
- Do not call Gemini or any LLM directly from tools — always use llm.Provider interface

## Project layout
```
cmd/server/
└── main.go                      ← entrypoint, wires everything together

internal/
├── llm/
│   ├── llm.go                   ← Provider interface
│   ├── gemini.go                ← GeminiProvider implementation
│   ├── config.go                ← FromConfig() reads LLM_PROVIDER/LLM_API_KEY/LLM_MODEL
│   └── llm_test.go              ← config tests + MockProvider for tool tests
├── prompts/
│   ├── prompts.go               ← go:embed directives + render functions
│   ├── evaluate_job.system.txt  ← static system prompt for evaluate_job
│   ├── evaluate_job.user.tmpl   ← dynamic user prompt template
│   ├── cover_letter.system.tmpl ← dynamic system prompt (tone varies)
│   ├── cover_letter.user.tmpl   ← dynamic user prompt template
│   └── prompts_test.go          ← template rendering tests
├── profile/
│   └── profile.go               ← Profile struct + Loader (reads data/profile.json)
├── job/
│   └── job.go                   ← Job + Compensation structs
├── tools/
│   ├── profile.go               ← get_candidate_profile handler
│   ├── util.go                  ← stripJSONFences() shared helper
│   ├── evaluate.go              ← evaluate_job tool + handler
│   └── evaluate_test.go         ← tests using MockProvider
└── db/
    └── (stub — do not implement until ticket 86c9wt0uw)

data/
├── profile.json                 ← real profile (gitignored)
└── profile.example.json         ← anonymized template (committed)
```

## Environment variables
```bash
LLM_PROVIDER=gemini              # supported: gemini
LLM_API_KEY=your-key-here        # from aistudio.google.com
LLM_MODEL=gemini-2.5-flash       # default if omitted
```

## Running locally
```bash
go run ./cmd/server
```

## Testing
```bash
go test ./...
```

## Connecting to Claude Desktop
Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "call-my-agent": {
      "command": "go",
      "args": ["run", "/absolute/path/to/call-my-agent/cmd/server"]
    }
  }
}
```

## Current state

### Done ✅
- internal/llm/      — Provider interface + Gemini + config + tests + MockProvider
- internal/prompts/  — go:embed + render functions + 4 prompt files + tests
- internal/profile/  — Profile struct + Loader
- internal/job/      — Job + Compensation structs
- cmd/server/        — MCP server skeleton, get_candidate_profile registered

### In progress 🔄
- internal/tools/evaluate.go    — evaluate_job tool (ticket 86c9wt0gj)
  - util.go             ✅ done (stripJSONFences)
  - evaluate.go         ✅ done (tool + handler)
  - evaluate_test.go    ⬜ next step
  - main.go wiring      ⬜ pending

### Not started ⬜
- internal/tools/coverletter.go — generate_cover_letter (ticket 86c9wt0nj)
- internal/db/                  — SQLite tracking (ticket 86c9wt0uw)
- .github/workflows/ci.yml      — GitHub Actions CI (ticket 86c9xj877)

## Tool registry
| Tool                  | Status      | Ticket    |
|-----------------------|-------------|-----------|
| get_candidate_profile | ✅ done     | 86c9wt05a |
| evaluate_job          | 🔄 in prog  | 86c9wt0gj |
| generate_cover_letter | ⬜ todo     | 86c9wt0nj |
| save_job              | ⬜ todo     | 86c9wt0uw |
| list_jobs             | ⬜ todo     | 86c9wt0uw |

## ClickUp project
Space: blackmagicbox
Folder: call-my-agent
List ID: 901523458793

## Slack
Channel: #call-my-agent-releases (GitHub notifications)
Channel: #call-my-agent (general)

## Instructions for Claude Code
When starting a new session, before doing anything else:

1. Run `find . -type f -name "*.go" | sort` to see what files exist
2. Compare the output against the "Current state" section above
3. For any file marked ✅ done — read it before touching it
4. For any file marked 🔄 in progress — read it, understand what's there, identify what's missing
5. For any file marked ⬜ todo — do not create it unless explicitly asked
6. Report back: "Here is what I found vs what CLAUDE.md expects" before writing any code
7. Never assume the state in CLAUDE.md is accurate — always verify against actual files

## Key context
- Dual purpose project: real job hunting tool + Google SRE L5/L6 portfolio artifact
- Target roles: Senior SRE, L5/L6, Berlin/Munich/Remote
- Next Google interview cycle: January 2027
- The scaffold pattern from this project will be reused across future projects
- MockProvider for tests lives in internal/llm/llm_test.go — import it in tool tests
